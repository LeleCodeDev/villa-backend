package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

type BookingService struct {
	txManager        *repository.TxManager
	repo             *repository.BookingRepository
	bookingAddonRepo *repository.BookingAddonRepository
	rentalOptionRepo *repository.RentalOptionRepository
	reasonRepo       *repository.ReasonRepository
	rentalPriceRepo  *repository.RentalPriceRepository
	seasonPeriodRepo *repository.SeasonPeriodRepository
	addonRepo        *repository.AddonRepository
	paymentRepo      *repository.PaymentRepository
	gateaway         *payment.MidtransGateway
}

func NewBookingService(
	txManager *repository.TxManager,
	repo *repository.BookingRepository,
	bookingAddonRepo *repository.BookingAddonRepository,
	rentalOptionRepo *repository.RentalOptionRepository,
	reasonRepo *repository.ReasonRepository,
	rentalPriceRepo *repository.RentalPriceRepository,
	seasonPeriodRepo *repository.SeasonPeriodRepository,
	addonRepo *repository.AddonRepository,
	paymentRepo *repository.PaymentRepository,
	gateaway *payment.MidtransGateway,
) *BookingService {
	return &BookingService{
		txManager:        txManager,
		repo:             repo,
		bookingAddonRepo: bookingAddonRepo,
		rentalOptionRepo: rentalOptionRepo,
		reasonRepo:       reasonRepo,
		rentalPriceRepo:  rentalPriceRepo,
		seasonPeriodRepo: seasonPeriodRepo,
		addonRepo:        addonRepo,
		paymentRepo:      paymentRepo,
		gateaway:         gateaway,
	}
}

func (s *BookingService) GetAll(ctx context.Context, query dto.BookingQuery, userID *uint) ([]dto.BookingResponse, int64, error) {
	bookings, total, err := s.repo.GetAll(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	bookingIds := make([]uint, 0, len(bookings))
	for _, booking := range bookings {
		bookingIds = append(bookingIds, booking.ID)
	}

	bookingAddons, err := s.bookingAddonRepo.GetAllInBookingIds(ctx, bookingIds)
	if err != nil {
		return nil, 0, err
	}

	payments, err := s.paymentRepo.GetAllInBookingIDs(ctx, bookingIds)
	if err != nil {
		return nil, 0, err
	}

	bookingAddonMap := make(map[uint][]model.BookingAddon, len(bookings))
	for _, bookingAddon := range bookingAddons {
		bookingAddonMap[bookingAddon.BookingID] = append(bookingAddonMap[bookingAddon.BookingID], bookingAddon)
	}

	paymentMap := make(map[uint][]model.Payment, len(payments))
	for _, payment := range payments {
		paymentMap[*payment.BookingID] = append(paymentMap[*payment.BookingID], payment)
	}

	responses := make([]dto.BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		responses = append(responses, mapper.ToBookingResponse(
			&booking,
			bookingAddonMap[booking.ID],
			paymentMap[booking.ID],
		))
	}

	return responses, total, nil
}

func (s *BookingService) GetByID(ctx context.Context, id uint, currentUser model.User) (dto.BookingResponse, error) {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.BookingResponse{}, err
	}
	if booking == nil || (currentUser.Role != model.RoleAdmin && booking.UserID != currentUser.ID) {
		return dto.BookingResponse{}, appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
	}

	bookingAddons, err := s.bookingAddonRepo.GetAllByBookingID(ctx, id)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	payments, err := s.paymentRepo.GetAllByBookingID(ctx, id)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(booking, bookingAddons, payments), nil
}

func parseBookingDates(in, out string) (time.Time, time.Time, error) {
	loc := config.Env.Location

	checkIn, err := time.ParseInLocation("2006-01-02", in, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid check_in_date, expected YYYY-MM-DD")
	}

	checkOut, err := time.ParseInLocation("2006-01-02", out, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid check_out_date, expected YYYY-MM-DD")
	}

	if checkOut.Before(checkIn) {
		return time.Time{}, time.Time{}, appError.BadRequest("check_out_date must be after check_in_date")
	}

	now := time.Now().In(config.Env.Location)
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, config.Env.Location)

	if checkIn.Before(today) {
		return time.Time{}, time.Time{}, appError.BadRequest("check_in_date cannot be in the past")
	}

	return checkIn, checkOut, nil
}

type priceKey struct {
	Season  model.SeasonType
	DayType model.DayType
}

func calculateStay(prices []model.RentalPrice, seasonPeriods []model.SeasonPeriod, checkIn, checkOut time.Time) (int64, error) {
	priceMap := make(map[priceKey]int64, len(prices))
	for _, r := range prices {
		priceMap[priceKey{Season: r.Season, DayType: r.DayType}] = r.Price
	}

	var total int64
	for d := checkIn; !d.After(checkOut); d = d.AddDate(0, 0, 1) {

		season := model.SeasonReguler
		for _, sp := range seasonPeriods {
			if !d.Before(sp.StartDate) && !d.After(sp.EndDate) {
				season = model.SeasonHigh
				break
			}
		}

		dayType := model.DayWeekday
		switch d.Weekday() {
		case time.Friday, time.Saturday, time.Sunday:
			dayType = model.DayWeekend
		}

		price, ok := priceMap[priceKey{Season: season, DayType: dayType}]
		if !ok {
			return 0, appError.BadRequest(fmt.Sprintf("no %s %s price configured for this rental option", season, dayType))
		}
		total += price
	}
	return total, nil
}

func (s *BookingService) Create(ctx context.Context, currentUser model.User, req dto.BookingCreateRequest) (dto.BookingResponse, error) {
	// parse req time from string to time.Time
	checkIn, checkOut, err := parseBookingDates(req.CheckInDate, req.CheckOutDate)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	// check addon duplicate
	seen := make(map[uint]struct{}, len(req.BookingAddons))
	addonIDs := make([]uint, 0, len(req.BookingAddons))
	for _, a := range req.BookingAddons {
		if _, dup := seen[a.AddonID]; dup {
			return dto.BookingResponse{}, appError.BadRequest("duplicate addon in booking_addons")
		}
		seen[a.AddonID] = struct{}{}
		addonIDs = append(addonIDs, a.AddonID)
	}

	const paymentWindow = 60 * time.Minute

	var createdBooking *model.Booking
	var createdAddons []model.BookingAddon
	var createdPayments []model.Payment

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txBookingRepo := s.repo.WithTx(tx)
		txBookingAddonRepo := s.bookingAddonRepo.WithTx(tx)
		txRentalOptionRepo := s.rentalOptionRepo.WithTx(tx)
		txReasonRepo := s.reasonRepo.WithTx(tx)
		txRentalPriceRepo := s.rentalPriceRepo.WithTx(tx)
		txSeasonPeriodRepo := s.seasonPeriodRepo.WithTx(tx)
		txAddonRepo := s.addonRepo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		// get rental option
		rentalOption, err := txRentalOptionRepo.GetByID(ctx, req.RentalOptionID)
		if err != nil {
			return err
		}
		if rentalOption == nil {
			return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", req.RentalOptionID))
		}

		if rentalOption.MaxCapacity < req.GuestCount {
			return appError.BadRequest(fmt.Sprintf("Max capacity is %d guest", rentalOption.MaxCapacity))
		}

		// get reason
		reason, err := txReasonRepo.GetByID(ctx, req.ReasonID)
		if err != nil {
			return err
		}
		if reason == nil {
			return appError.NotFound(fmt.Sprintf("Reason not found with ID: %d", req.ReasonID))
		}

		// check booking conflict
		conflicts, err := txBookingRepo.GetBookingConflicts(ctx, rentalOption.ID, checkIn, checkOut, 0)
		if err != nil {
			return err
		}
		if len(conflicts) > 0 {
			conflictResponses := make([]dto.BookingConflict, 0, len(conflicts))
			for _, conflict := range conflicts {
				conflictResponses = append(conflictResponses, dto.BookingConflict{CheckInDate: conflict.CheckInDate, CheckOutDate: conflict.CheckOutDate})
			}

			return appError.Conflict("rental option is already booked for these dates").
				WithErrors(conflictResponses)
		}

		// get rental prices
		prices, err := txRentalPriceRepo.GetAllByRentalOptionID(ctx, rentalOption.ID)
		if err != nil {
			return err
		}

		// get season periods
		seasonPeriods, err := txSeasonPeriodRepo.GetSeasonConflict(ctx, checkIn, checkOut, 0)
		if err != nil {
			return err
		}

		// calculate booking price
		bookingPrice, err := calculateStay(prices, seasonPeriods, checkIn, checkOut)
		if err != nil {
			return err
		}

		// get addons and calculate addon price
		addons, err := txAddonRepo.GetAllInIds(ctx, addonIDs)
		if err != nil {
			return err
		}
		addonMap := make(map[uint]model.Addon, len(addons))
		for _, a := range addons {
			addonMap[a.ID] = a
		}

		var addonTotal int64
		for _, ba := range req.BookingAddons {
			a, ok := addonMap[ba.AddonID]
			if !ok {
				return appError.NotFound(fmt.Sprintf("addon not found with id: %d", ba.AddonID))
			}
			addonTotal += a.Price * int64(ba.Quantity)
		}

		totalPrice := bookingPrice + addonTotal

		items := []midtrans.ItemDetails{{
			ID:    fmt.Sprintf("STAY-%d", rentalOption.ID),
			Name:  rentalOption.Title,
			Price: bookingPrice,
			Qty:   1,
		}}

		for _, ba := range req.BookingAddons {
			a := addonMap[ba.AddonID]
			items = append(items, midtrans.ItemDetails{
				ID:    fmt.Sprintf("ADDON-%d", a.ID),
				Name:  a.Title,
				Price: a.Price,
				Qty:   int32(ba.Quantity),
			})
		}

		// now := time.Now()
		payKind := model.PaymentKindFull
		payAmount := totalPrice
		var depositAmount *int64
		var balanceDueAt *time.Time

		// check if its a deposit
		if req.PaymentType == model.PaymentTypeDeposit {
			d := totalPrice * 10 / 100
			depositAmount = &d

			payKind = model.PaymentKindDeposit
			payAmount = d

			bda := checkIn.Add(11 * time.Hour)
			balanceDueAt = &bda

			items = []midtrans.ItemDetails{{
				ID:    fmt.Sprintf("DEPOSIT-BOOKING-%d", rentalOption.ID),
				Name:  "Deposit 10% - " + rentalOption.Title,
				Price: d,
				Qty:   1,
			}}
		}

		// create booking
		booking := mapper.ToBookingModel(
			req,
			currentUser,
			*rentalOption,
			*reason,
			checkIn,
			checkOut,
			bookingPrice,
			bookingPrice+addonTotal,
			depositAmount,
			balanceDueAt,
			time.Now().Add(paymentWindow),
		)
		if err := txBookingRepo.Create(ctx, booking); err != nil {
			return err
		}

		// create booking addons
		bookingAddons := make([]model.BookingAddon, 0, len(req.BookingAddons))
		if len(req.BookingAddons) > 0 {
			for _, ar := range req.BookingAddons {
				ba := mapper.ToBookingAddonModel(ar, *booking, addonMap[ar.AddonID])
				bookingAddons = append(bookingAddons, *ba)
			}
			if err := txBookingAddonRepo.Creates(ctx, bookingAddons); err != nil {
				return err
			}
		}

		// create payment
		orderID := fmt.Sprintf("BOOKING-%d-%s-%d", booking.ID, strings.ToUpper(string(payKind)), time.Now().UnixNano())
		token, redirectURL, err := s.gateaway.CreateSnapToken(
			orderID,
			payAmount,
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
			OrderID:     orderID,
			BookingID:   &booking.ID,
			Kind:        payKind,
			Status:      model.PaymentStatusPending,
			Amount:      payAmount,
			SnapToken:   &token,
			RedirectURL: &redirectURL,
			ExpiresAt:   booking.ExpiresAt,
		}
		if err := txPaymentRepo.Create(ctx, payment); err != nil {
			return err
		}

		createdBooking = booking
		createdAddons = bookingAddons
		createdPayments = append(createdPayments, *payment)

		return nil
	}); err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(createdBooking, createdAddons, createdPayments), nil
}

func (s *BookingService) PayBalance(ctx context.Context, currentUser model.User, id uint) (dto.PaymentResponse, error) {
	var balancePayment *model.Payment

	const paymentWindow = 60 * time.Minute

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txBookingRepo := s.repo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		// get booking
		booking, err := txBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if booking == nil || booking.UserID != currentUser.ID {
			return appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
		}
		if booking.Status != model.BookingStatusPartiallyPaid {
			return appError.BadRequest("booking is not awaiting balance payment")
		}
		if booking.DepositAmount == nil {
			return appError.BadRequest("booking has no deposit")
		}

		// check existing pending balance payment
		existing, err := txPaymentRepo.GetFirstByBookingIDAndStatusAndKind(ctx, booking.ID, model.PaymentStatusPending, model.PaymentKindBalance)
		if err != nil {
			return err
		}
		if existing != nil && existing.ExpiresAt != nil && existing.ExpiresAt.After(time.Now()) {
			balancePayment = existing
			return nil
		}

		// make midtrans item detail
		itemID := fmt.Sprintf("BALANCE-BOOKING-#%d", booking.ID)
		itemName := fmt.Sprintf("Balance #%d", booking.ID)

		if booking.RentalOptionID != nil && booking.RentalOption != nil {
			itemID = fmt.Sprintf("BALANCE-BOOKING-%d", booking.RentalOptionID)
			itemName = "Balance - " + booking.RentalOption.Title
		}

		balance := booking.TotalPrice - *booking.DepositAmount

		items := []midtrans.ItemDetails{{
			ID:    itemID,
			Name:  itemName,
			Price: balance,
			Qty:   1,
		}}

		// make snap token
		orderID := fmt.Sprintf("BOOKING-%d-BALANCE-%d", booking.ID, time.Now().UnixNano())
		token, redirectURL, err := s.gateaway.CreateSnapToken(
			orderID,
			balance,
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

		// make payment model
		expiresAt := time.Now().Add(paymentWindow)
		payment := &model.Payment{
			OrderID:     orderID,
			BookingID:   &booking.ID,
			Kind:        model.PaymentKindBalance,
			Status:      model.PaymentStatusPending,
			Amount:      balance,
			SnapToken:   &token,
			RedirectURL: &redirectURL,
			ExpiresAt:   &expiresAt,
		}
		if err := txPaymentRepo.Create(ctx, payment); err != nil {
			return err
		}

		balancePayment = payment

		return nil
	}); err != nil {
		return dto.PaymentResponse{}, err
	}

	return mapper.ToPaymentResponse(balancePayment), nil
}

func (s *BookingService) Cancel(ctx context.Context, currentUser model.User, id uint) (dto.BookingResponse, error) {
	var canceledBooking *model.Booking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txBookingRepo := s.repo.WithTx(tx)
		txPaymentRepo := s.paymentRepo.WithTx(tx)

		booking, err := txBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if booking == nil || (currentUser.Role != model.RoleAdmin && booking.UserID != currentUser.ID) {
			return appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
		}

		if booking.Status != model.BookingStatusPending {
			return appError.BadRequest("Booking status is not pending!")
		}

		moved, err := txBookingRepo.TransitionStatus(ctx, booking.ID, model.BookingStatusPending, model.BookingStatusCanceled)
		if err != nil {
			return err
		}
		if !moved {
			return appError.Conflict("Booking status has changed!")
		}

		pendingPayments, err := txPaymentRepo.GetAllByBookingIDAndStatus(ctx, booking.ID, model.PaymentStatusPending)
		if err != nil {
			return err
		}

		for _, pendingPayment := range pendingPayments {
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

		booking.Status = model.BookingStatusCanceled
		canceledBooking = booking
		return nil
	}); err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(canceledBooking, nil, nil), nil
}

func (s *BookingService) CheckIn(ctx context.Context, id uint) (dto.BookingResponse, error) {
	var paidBooking *model.Booking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txBookingRepo := s.repo.WithTx(tx)

		booking, err := txBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
		}

		if booking.Status != model.BookingStatusPaid {
			return appError.BadRequest("Booking status is not paid")
		}

		now := time.Now()
		booking.CheckInAt = &now
		if err := txBookingRepo.Update(ctx, booking); err != nil {
			return err
		}

		paidBooking = booking

		return nil
	}); err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(paidBooking, nil, nil), nil
}

func (s *BookingService) Complete(ctx context.Context, id uint) (dto.BookingResponse, error) {
	var completeBooking *model.Booking

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txBookingRepo := s.repo.WithTx(tx)

		booking, err := txBookingRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
		}

		if booking.Status != model.BookingStatusPaid {
			return appError.BadRequest("Booking status is not paid")
		}

		now := time.Now()
		booking.CheckOutAt = &now
		booking.Status = model.BookingStatusCompleted
		if err := txBookingRepo.Update(ctx, booking); err != nil {
			return err
		}

		completeBooking = booking

		return nil
	}); err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(completeBooking, nil, nil), nil
}
