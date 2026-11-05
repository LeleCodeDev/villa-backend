package service

import (
	"context"
	"errors"
	"log"
	"math"
	"strconv"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAmountMismatch   = errors.New("amount mismatch")
)

type PaymentService struct {
	repo    *repository.BookingRepository
	gateway *payment.MidtransGateway
}

func NewPaymentService(
	repo *repository.BookingRepository,
	gateway *payment.MidtransGateway,
) *PaymentService {
	return &PaymentService{
		repo:    repo,
		gateway: gateway,
	}
}

func (s *PaymentService) HandleNotification(ctx context.Context, n dto.MidtransNotification) error {
	booking, err := s.repo.GetByOrderID(ctx, n.OrderID)
	if err != nil {
		return err
	}
	if booking == nil {
		return ErrBookingNotFound
	}

	status, err := s.gateway.CheckTransaction(ctx, n.OrderID)
	if err != nil {
		return err
	}

	gross, err := strconv.ParseFloat(n.GrossAmount, 64)
	if err != nil || int64(math.Round(gross)) != booking.TotalPrice {
		return ErrAmountMismatch
	}

	target, ok := mapStatus(status.TransactionStatus, status.FraudStatus)
	if !ok || booking.Status == target {
		return nil
	}

	updated, err := s.repo.TransitionStatus(ctx, booking.ID, model.BookingStatusPending, target)
	if err != nil {
		return err
	}
	if !updated {
		if target == model.BookingStatusPaid {
			log.Printf("REVIEW NEEDED: payment for booking %d arrived while status is %s", booking.ID, booking.Status)
		}
		return nil
	}

	// if target == model.BookingStatusPaid {
	// }
	return nil
}

func mapStatus(transactionStatus, fraudStatus string) (model.BookingStatus, bool) {
	switch transactionStatus {
	case "settlement":
		return model.BookingStatusPaid, true
	case "capture":
		if fraudStatus == "accept" {
			return model.BookingStatusPaid, true
		}
	case "expire":
		return model.BookingStatusExpired, true
	case "cancel":
		return model.BookingStatusCanceled, true
	}
	return "", false
}
