package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/midtrans/midtrans-go"
	"gorm.io/gorm"
)

type EventBookingService struct {
	txManager        *repository.TxManager
	repo             *repository.EventBookingRepository
	villaPackageRepo *repository.VillaPackageRepository
	paymentRepo      *repository.PaymentRepository
	gateaway         *payment.MidtransGateway
}

func NewEventBookingService(
	txManager *repository.TxManager,
	repo *repository.EventBookingRepository,
	villaPackageRepo *repository.VillaPackageRepository,
	paymentRepo *repository.PaymentRepository,
	gateaway *payment.MidtransGateway,
) *EventBookingService {
	return &EventBookingService{
		txManager:        txManager,
		repo:             repo,
		villaPackageRepo: villaPackageRepo,
		paymentRepo:      paymentRepo,
		gateaway:         gateaway,
	}
}

func (s *EventBookingService) GetAll(ctx context.Context, query dto.EventBookingQuery, userID *uint) ([]dto.EventBookingResponse, int64, error) {
	eventBokings, total, err := s.repo.GetAll(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	eventBookingIds := make([]uint, 0, len(eventBokings))

	payments, err := s.paymentRepo.GetAllInEventBookingIds(ctx, eventBookingIds)
	if err != nil {
		return nil, 0, err
	}

	paymentMap := make(map[uint]*model.Payment, len(eventBokings))
	for _, payment := range payments {
		paymentMap[*payment.EventBookingID] = &payment
	}

	responses := make([]dto.EventBookingResponse, 0, len(eventBokings))
	for _, eventBoking := range eventBokings {
		responses = append(responses, mapper.ToEventBookingResponse(&eventBoking, paymentMap[eventBoking.ID]))
	}

	return responses, total, nil
}

func (s *EventBookingService) GetByID(ctx context.Context, id uint, currentUser model.User) (dto.EventBookingResponse, error) {
	eventBooking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.EventBookingResponse{}, err
	}
	if eventBooking == nil {
		return dto.EventBookingResponse{}, appError.NotFound(fmt.Sprintf("Event booking not found with ID : %d", id))
	}

	if currentUser.Role != model.RoleAdmin && eventBooking.UserID != currentUser.ID {
		return dto.EventBookingResponse{}, appError.Unauthorized("You are not authorized to access this event booking!")
	}

	payment, err := s.paymentRepo.GetByEventBookingID(ctx, eventBooking.ID)
	if err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(eventBooking, payment), nil
}

func parseEventBookingDates(start, end string) (time.Time, time.Time, error) {
	loc := config.Env.Location

	startAt, err := time.ParseInLocation("2006-01-02 15:04", start, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid start_at, expected YYYY-MM-DD HH:MM")
	}

	endAt, err := time.ParseInLocation("2006-01-02 15:04", end, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid end_at, expected YYYY-MM-DD HH:MM")
	}

	if endAt.Before(startAt) {
		return time.Time{}, time.Time{}, appError.BadRequest("end_at must be after start_at")
	}

	now := time.Now().In(config.Env.Location)
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, config.Env.Location)

	if startAt.Before(today) {
		return time.Time{}, time.Time{}, appError.BadRequest("start_at cannot be in the past")
	}

	return startAt, endAt, nil
}

func calculateQuantity(t model.ChargeType, start, end time.Time) int {
	d := end.Sub(start)
	var q int
	switch t {
	case model.ChargePerHour:
		q = int(math.Ceil(d.Hours()))
	case model.ChargePerDay:
		q = int(math.Ceil(d.Hours() / 24))
	default:
		q = 1
	}
	if q < 1 {
		q = 1
	}
	return q
}

func (s *EventBookingService) Create(ctx context.Context, currentUser model.User, req dto.EventBookingCreateRequest) (dto.EventBookingResponse, error) {
	// parse req string time from string to time.Time
	startAt, endAt, err := parseEventBookingDates(req.StartAt, req.EndAt)
	if err != nil {
		return dto.EventBookingResponse{}, err
	}

	const paymentWindow = 60 * time.Minute

	var createdEventBooking *model.EventBooking
	var createdPayment *model.Payment

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txEventBookingRepo := s.repo.WithTx(tx)
		txVillaPackageRepo := s.villaPackageRepo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		// get villa package
		villaPackage, err := txVillaPackageRepo.GetByID(ctx, req.VillaPackageID)
		if err != nil {
			return err
		}
		if villaPackage == nil {
			return appError.NotFound(fmt.Sprintf("Villa package not found with ID: %d", req.VillaPackageID))
		}

		if villaPackage.RequiresApproval {
			return appError.BadRequest("Villa package requires approval, please submit a request")
		}

		// check booking conflict
		conflicts, err := txEventBookingRepo.GetBookingConflicts(ctx, startAt, endAt, 0)
		if err != nil {
			return err
		}
		if len(conflicts) > 0 {
			conflictResponses := make([]dto.EventBookingConflict, 0, len(conflicts))
			for _, conflict := range conflicts {
				conflictResponses = append(conflictResponses, dto.EventBookingConflict{StartAt: conflict.StartAt, EndAt: conflict.EndAt})
			}

			return appError.Conflict("Event is already booked for these dates").
				WithErrors(conflictResponses)
		}

		// calculate price
		quantity := calculateQuantity(villaPackage.Type, startAt, endAt)
		totalPrice := villaPackage.Price * int64(quantity)

		items := []midtrans.ItemDetails{{
			ID:    fmt.Sprintf("EVENT-%s-%d", villaPackage.Type, villaPackage.ID),
			Name:  villaPackage.Title,
			Price: villaPackage.Price,
			Qty:   int32(quantity),
		}}

		// create event booking
		eventBooking := mapper.ToEventBookingModel(
			startAt,
			endAt,
			villaPackage.Title,
			currentUser,
			villaPackage,
			req.GuestCount,
			quantity,
			villaPackage.Price,
			totalPrice,
			time.Now().Add(paymentWindow),
		)
		if err := txEventBookingRepo.Create(ctx, eventBooking); err != nil {
			return err
		}

		// create payment
		orderID := fmt.Sprintf("EVENT-BOOKING-%d-%d", eventBooking.ID, time.Now().UnixNano())
		token, redirectURL, err := s.gateaway.CreateSnapToken(
			orderID,
			eventBooking.TotalPrice,
			items,
			&midtrans.CustomerDetails{
				FName: currentUser.Username,
				Email: currentUser.Email,
				Phone: currentUser.Phone,
			},
			int(paymentWindow.Minutes()))
		if err != nil {
			return appError.BadGateway("failed to start payment, please try again")
		}

		payment := &model.Payment{
			OrderID:        orderID,
			EventBookingID: &eventBooking.ID,
			Kind:           model.PaymentKindFull,
			Status:         model.PaymentStatusPending,
			Amount:         totalPrice,
			SnapToken:      &token,
			RedirectURL:    &redirectURL,
			ExpiresAt:      eventBooking.ExpiresAt,
		}
		if err := txPaymentRepo.Create(ctx, payment); err != nil {
			return err
		}

		createdEventBooking = eventBooking
		createdPayment = payment
		return nil
	}); err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(createdEventBooking, createdPayment), nil
}

func (s *EventBookingService) Cancel(ctx context.Context, currentUser model.User, id uint) (dto.EventBookingResponse, error) {
	var canceledBooking *model.EventBooking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txEventBookingRepo := s.repo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		eventBooking, err := txEventBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if eventBooking == nil || (currentUser.Role != model.RoleAdmin && eventBooking.UserID != currentUser.ID) {
			return appError.NotFound(fmt.Sprintf("Event booking not found with ID: %d", id))
		}

		if eventBooking.Status != model.EventBookingStatusPending {
			return appError.BadRequest("Event booking status is not pending!")
		}

		moved, err := txEventBookingRepo.TransitionStatus(ctx, eventBooking.ID, model.EventBookingStatusPending, model.EventBookingStatusCanceled)
		if err != nil {
			return err
		}
		if !moved {
			return appError.Conflict("Booking status has changed!")
		}

		pendingPayment, err := txPaymentRepo.GetByEventBookingIDAndStatus(ctx, eventBooking.ID, model.PaymentStatusPending)
		if err != nil {
			return err
		}
		if pendingPayment != nil {
			if err := s.gateaway.CancelTransaction(ctx, pendingPayment.OrderID); err != nil {
				if errors.Is(err, payment.ErrNotCancelable) {
					return appError.Conflict("Payment was already received, booking can't be canceled!")
				}

				return appError.BadGateway("Failed to cancel payment, please try again")
			}

			if _, err := txPaymentRepo.TransitionStatus(
				ctx,
				pendingPayment.ID,
				model.PaymentStatusPending,
				model.PaymentStatusCanceled); err != nil {
				return err
			}
		}

		eventBooking.Status = model.EventBookingStatusCanceled
		canceledBooking = eventBooking

		return nil
	}); err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(canceledBooking, nil), nil
}

func (s *EventBookingService) CheckIn(ctx context.Context, id uint) (dto.EventBookingResponse, error) {
	var paidBooking *model.EventBooking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txEventBookingRepo := s.repo.WithTx(tx)

		eventBooking, err := txEventBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if eventBooking == nil {
			return appError.NotFound(fmt.Sprintf("Event booking not found with ID: %d", id))
		}

		if eventBooking.Status != model.EventBookingStatusPaid {
			return appError.BadRequest("Event booking status is not paid")
		}

		now := time.Now()
		eventBooking.CheckInAt = &now
		if err := txEventBookingRepo.Update(ctx, eventBooking); err != nil {
			return err
		}

		paidBooking = eventBooking

		return nil
	}); err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(paidBooking, nil), nil
}

func (s *EventBookingService) Complete(ctx context.Context, id uint) (dto.EventBookingResponse, error) {
	var completeBooking *model.EventBooking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txEventBookingRepo := s.repo.WithTx(tx)

		eventBooking, err := txEventBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if eventBooking == nil {
			return appError.NotFound(fmt.Sprintf("Event booking not found with ID: %d", id))
		}

		if eventBooking.Status != model.EventBookingStatusPaid {
			return appError.BadRequest("Booking status is not paid")
		}

		now := time.Now()
		eventBooking.CheckOutAt = &now
		eventBooking.Status = model.EventBookingStatusCompleted
		if err := txEventBookingRepo.Update(ctx, eventBooking); err != nil {
			return err
		}

		completeBooking = eventBooking

		return nil
	}); err != nil {
		return dto.EventBookingResponse{}, err
	}

	return mapper.ToEventBookingResponse(completeBooking, nil), nil
}
