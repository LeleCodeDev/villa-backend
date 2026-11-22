package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type EventBookingRequestService struct {
	txManager        *repository.TxManager
	repo             *repository.EventBookingRequestRepository
	villaPackageRepo *repository.VillaPackageRepository
	eventBookingRepo *repository.EventBookingRepository
	paymentRepo      *repository.PaymentRepository
	gateaway         *payment.MidtransGateway
}

func NewEventBookingRequestService(
	txManager *repository.TxManager,
	repo *repository.EventBookingRequestRepository,
	villaPackageRepo *repository.VillaPackageRepository,
	eventBookingRepo *repository.EventBookingRepository,
	paymentRepo *repository.PaymentRepository,
	gateaway *payment.MidtransGateway,
) *EventBookingRequestService {
	return &EventBookingRequestService{
		txManager:        txManager,
		repo:             repo,
		villaPackageRepo: villaPackageRepo,
		eventBookingRepo: eventBookingRepo,
		paymentRepo:      paymentRepo,
		gateaway:         gateaway,
	}
}

const (
	paymentWindow = 60 * time.Minute
	packageHold   = 24 * time.Hour
	customHold    = 48 * time.Hour
	minLeadTime   = 4 * 24 * time.Hour
)

func (s *EventBookingRequestService) GetAll(ctx context.Context, query dto.EventBookingRequestQuery, userID *uint) ([]dto.EventBookingRequestResponse, int64, error) {
	requests, total, err := s.repo.GetAll(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.EventBookingRequestResponse, 0, len(requests))
	for _, request := range requests {
		responses = append(responses, mapper.ToEventBookingRequestResponse(&request))
	}

	return responses, total, nil
}

func (s *EventBookingRequestService) GetByID(ctx context.Context, id uint, currentUser model.User) (dto.EventBookingRequestResponse, error) {
	request, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.EventBookingRequestResponse{}, err
	}
	if request == nil || (currentUser.Role != model.RoleAdmin && request.UserID != currentUser.ID) {
		return dto.EventBookingRequestResponse{}, appError.NotFound(fmt.Sprintf("Event request not found with ID: %d", id))
	}

	return mapper.ToEventBookingRequestResponse(request), nil
}

func (s *EventBookingRequestService) Create(ctx context.Context, req dto.CreateEventRequest, currentUser model.User) (dto.EventBookingRequestResponse, error) {
	const minLeadTime = 4 * 24 * time.Hour

	startAt, endAt, err := parseEventBookingDates(req.StartAt, req.EndAt)
	if err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	now := time.Now()
	if startAt.Sub(now) < minLeadTime {
		return dto.EventBookingRequestResponse{}, appError.BadRequest(
			fmt.Sprintf("events must be requested at least %d days in advance", int(minLeadTime.Hours()/24)))
	}

	var createdRequest *model.EventBookingRequest

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txVillaPackageRepo := s.villaPackageRepo.WithTx(tx)

		var villaPackage *model.VillaPackage
		if req.VillaPackageID != nil {
			villaPackage, err = txVillaPackageRepo.GetByID(ctx, *req.VillaPackageID)
			if err != nil {
				return err
			}
			if villaPackage == nil {
				return appError.BadRequest(fmt.Sprintf("Villa package not found with ID : %d", *req.VillaPackageID))
			}
		}

		request := mapper.ToEventBookingRequestModel(req, currentUser, villaPackage, startAt, endAt)
		if err := txRepo.Create(ctx, request); err != nil {
			return err
		}

		createdRequest = request

		return nil
	}); err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	return mapper.ToEventBookingRequestResponse(createdRequest), nil
}

func (s *EventBookingRequestService) Approve(ctx context.Context, id uint, req dto.ApproveEventRequest) (dto.EventBookingResponse, error) {
	var approvedBooking *model.EventBooking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRequestRepo := s.repo.WithTx(tx)
		txEventBookingRepo := s.eventBookingRepo.WithTx(tx)

		request, err := txRequestRepo.GetByID(ctx, id) // preload VillaPackage + User, ideally FOR UPDATE
		if err != nil {
			return err
		}
		if request == nil {
			return appError.NotFound(fmt.Sprintf("Event booking request not found with ID : %d", id))
		}
		if request.Status != model.RequestPending {
			return appError.BadRequest("Event booking request status is not pending!")
		}

		now := time.Now()
		if request.StartAt.Sub(now) < minLeadTime {
			return appError.BadRequest("this request is too close to the event date to approve, please contact the customer")
		}

		var price int64
		quantity := 1
		hold := customHold

		if request.VillaPackageID != nil {
			villaPackage := request.VillaPackage
			if villaPackage == nil {
				return appError.BadRequest("villa package no longer exists")
			}
			if req.Price != nil {
				return appError.BadRequest("price cannot be set for a package event")
			}

			quantity = calculateQuantity(villaPackage.Type, request.StartAt, request.EndAt)
			price = villaPackage.Price
			hold = packageHold

			conflicts, err := txEventBookingRepo.GetBookingConflicts(ctx, villaPackage.ID, request.StartAt, request.EndAt, 0)
			if err != nil {
				return err
			}

			if len(conflicts) > 0 {
				conflictResponses := make([]dto.EventBookingConflict, 0, len(conflicts))
				for _, c := range conflicts {
					conflictResponses = append(conflictResponses, dto.EventBookingConflict{StartAt: c.StartAt, EndAt: c.EndAt})
				}

				return appError.Conflict("Villa package is already booked for these dates").WithErrors(conflictResponses)
			}
		} else {
			if req.Price == nil || *req.Price <= 0 {
				return appError.BadRequest("price is required for custom event!")
			}
			price = *req.Price

			conflicts, err := txEventBookingRepo.GetCustomBookingConflicts(ctx, request.StartAt, request.EndAt, 0)
			if err != nil {
				return err
			}

			if len(conflicts) > 0 {
				conflictResponses := make([]dto.EventBookingConflict, 0, len(conflicts))
				for _, c := range conflicts {
					conflictResponses = append(conflictResponses, dto.EventBookingConflict{StartAt: c.StartAt, EndAt: c.EndAt})
				}

				return appError.Conflict("Event is already booked for these dates").WithErrors(conflictResponses)
			}
		}

		expiresAt := now.Add(hold)
		if latest := request.StartAt.Add(-24 * time.Hour); expiresAt.After(latest) {
			expiresAt = latest
		}

		totalPrice := price * int64(quantity)

		eventBooking := mapper.ToEventBookingModel(
			request.StartAt,
			request.EndAt,
			request.Title,
			request.User,
			request.VillaPackage,
			request.GuestCount,
			quantity,
			price,
			totalPrice,
			expiresAt,
		)
		if err := txEventBookingRepo.Create(ctx, eventBooking); err != nil {
			return err
		}

		request.Status = model.RequestApproved
		request.ReviewedAt = &now
		request.EventBookingID = &eventBooking.ID
		if err := txRequestRepo.Update(ctx, request); err != nil {
			return err
		}

		approvedBooking = eventBooking
		return nil
	}); err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(approvedBooking, nil), nil
}

func (s *EventBookingRequestService) Pay(ctx context.Context, currentUser model.User, id uint) (dto.PaymentResponse, error) {
	var createdPayment *model.Payment

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRequestRepo := s.repo.WithTx(tx)
		txEventBookingRepo := s.eventBookingRepo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		request, err := txRequestRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if request == nil {
			return appError.NotFound(fmt.Sprintf("Event booking request not found with ID : %d", id))
		}
		if request.Status != model.RequestApproved {
			return appError.BadRequest("Event booking request status is not approved!")
		}

		eventBooking, err := txEventBookingRepo.GetByID(ctx, *request.EventBookingID)
		if err != nil {
			return err
		}
		if eventBooking == nil || eventBooking.UserID != currentUser.ID {
			return appError.NotFound(fmt.Sprintf("Event booking not found with ID: %d", *request.EventBookingID))
		}

		now := time.Now()
		if eventBooking.Status != model.EventBookingStatusPending || eventBooking.ExpiresAt == nil || !eventBooking.ExpiresAt.After(now) {
			return appError.BadRequest("Event booking is not payable")
		}

		existing, err := txPaymentRepo.GetByEventBookingIDAndStatus(ctx, eventBooking.ID, model.PaymentStatusPending)
		if err != nil {
			return err
		}
		if existing != nil && existing.ExpiresAt != nil && existing.ExpiresAt.After(now) {
			createdPayment = existing
			return nil
		}

		expiresAt := now.Add(paymentWindow)
		if expiresAt.After(*eventBooking.ExpiresAt) {
			expiresAt = *eventBooking.ExpiresAt
		}
		minutes := int(math.Ceil(expiresAt.Sub(now).Minutes()))
		if minutes < 1 {
			return appError.BadRequest("Event booking is not payable")
		}

		items := []midtrans.ItemDetails{{
			ID:    fmt.Sprintf("EVENT-%d", eventBooking.ID),
			Name:  fmt.Sprintf("Event Booking #%d", eventBooking.ID),
			Price: eventBooking.Price,
			Qty:   int32(eventBooking.Quantity),
		}}

		orderID := fmt.Sprintf("EVENT-BOOKING-%d-%d", eventBooking.ID, now.UnixNano())
		token, redirectURL, err := s.gateaway.CreateSnapToken(
			orderID,
			eventBooking.TotalPrice,
			items,
			&midtrans.CustomerDetails{
				FName: currentUser.Username,
				Email: currentUser.Email,
				Phone: currentUser.Phone,
			},
			minutes)
		if err != nil {
			return appError.BadGateway("failed to start payment, please try again")
		}

		payment := &model.Payment{
			OrderID:        orderID,
			EventBookingID: &eventBooking.ID,
			Kind:           model.PaymentKindFull,
			Status:         model.PaymentStatusPending,
			Amount:         eventBooking.TotalPrice,
			SnapToken:      &token,
			RedirectURL:    &redirectURL,
			ExpiresAt:      &expiresAt,
		}
		if err := txPaymentRepo.Create(ctx, payment); err != nil {
			return err
		}

		createdPayment = payment
		return nil
	}); err != nil {
		return dto.PaymentResponse{}, err
	}

	return mapper.ToPaymentResponse(createdPayment), nil
}

func (s *EventBookingRequestService) Reject(ctx context.Context, id uint, req dto.RejectEventRequest) (dto.EventBookingRequestResponse, error) {
	var rejected *model.EventBookingRequest

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		request, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if request == nil {
			return appError.NotFound(fmt.Sprintf("Event booking request not found with ID : %d", id))
		}
		if request.Status != model.RequestPending {
			return appError.BadRequest("Event booking request status is not pending!")
		}

		now := time.Now()
		reason := strings.TrimSpace(req.RejectReason)
		request.Status = model.RequestRejected
		request.ReviewedAt = &now
		request.RejectReason = &reason
		if err := txRepo.Update(ctx, request); err != nil {
			return err
		}

		rejected = request
		return nil
	}); err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	return mapper.ToEventBookingRequestResponse(rejected), nil
}

func (s *EventBookingRequestService) Cancel(ctx context.Context, currentUser model.User, id uint) (dto.EventBookingRequestResponse, error) {
	var canceled *model.EventBookingRequest

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		request, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if request == nil || request.UserID != currentUser.ID {
			return appError.NotFound(fmt.Sprintf("Event booking request not found with ID : %d", id))
		}
		if request.Status != model.RequestPending {
			return appError.BadRequest("Event booking request is not pending")
		}

		now := time.Now()
		request.Status = model.RequestCanceled
		request.ReviewedAt = &now
		if err := txRepo.Update(ctx, request); err != nil {
			return err
		}

		canceled = request
		return nil
	}); err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	return mapper.ToEventBookingRequestResponse(canceled), nil
}
