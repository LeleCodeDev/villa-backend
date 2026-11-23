package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type PaymentService struct {
	txManager        *repository.TxManager
	paymentRepo      *repository.PaymentRepository
	bookingRepo      *repository.BookingRepository
	eventBookingRepo *repository.EventBookingRepository
	gateway          *payment.MidtransGateway
}

func NewPaymentService(
	txManager *repository.TxManager,
	paymentRepo *repository.PaymentRepository,
	bookingRepo *repository.BookingRepository,
	eventBookingRepo *repository.EventBookingRepository,
	gateway *payment.MidtransGateway,
) *PaymentService {
	return &PaymentService{
		txManager:        txManager,
		paymentRepo:      paymentRepo,
		bookingRepo:      bookingRepo,
		eventBookingRepo: eventBookingRepo,
		gateway:          gateway,
	}
}

func (s *PaymentService) HandleNotificationCallback(ctx context.Context, n dto.MidtransCallbackRequest) error {
	p, err := s.paymentRepo.GetByOrderID(ctx, n.OrderID)
	if err != nil {
		return err
	}
	if p == nil {
		return appError.NotFound(fmt.Sprintf("Payment not found with order ID : %s", n.OrderID))
	}

	status, err := s.gateway.CheckTransaction(ctx, n.OrderID)
	if err != nil {
		return err
	}

	grossAmount, err := strconv.ParseFloat(n.GrossAmount, 64)
	if err != nil || int64(math.Round(grossAmount)) != p.Amount {
		return appError.Conflict("Amount mismatch")
	}

	target, ok := mapPaymentStatus(status.TransactionStatus, status.FraudStatus)
	if !ok || p.Status == target {
		return nil
	}

	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txPaymentRepo := s.paymentRepo.WithTx(tx)
		txBookingRepo := s.bookingRepo.WithTx(tx)
		txEventBookingRepo := s.eventBookingRepo.WithTx(tx)

		updated, err := txPaymentRepo.TransitionStatus(ctx, p.ID, model.PaymentStatusPending, target)
		if err != nil {
			return err
		}
		if !updated {
			if target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but was already %s", p.OrderID, p.Status)
			}

			return nil
		}

		if p.BookingID != nil {
			from, to, ok := bookingTransition(p.Kind, target)
			if !ok {
				return nil
			}

			moved, err := txBookingRepo.TransitionStatus(ctx, *p.BookingID, from, to)
			if err != nil {
				return err
			}

			if !moved && target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but booking not in status %s", p.OrderID, from)
			}
		}

		if p.EventBookingID != nil {
			from, to, ok := eventBookingTransition(p.Kind, target)
			if !ok {
				return nil
			}

			moved, err := txEventBookingRepo.TransitionStatus(ctx, *p.EventBookingID, from, to)
			if err != nil {
				return err
			}

			if !moved && target == model.PaymentStatusPaid {
				log.Printf("REVIEW NEEDED: payment %s paid but event booking not in status %s", p.OrderID, from)
			}
		}

		return nil
	})
}

func mapPaymentStatus(transactionStatus, fraudStatus string) (model.PaymentStatus, bool) {
	switch transactionStatus {
	case "settlement":
		return model.PaymentStatusPaid, true
	case "capture":
		if fraudStatus == "accept" {
			return model.PaymentStatusPaid, true
		}
	case "expire":
		return model.PaymentStatusExpired, true
	case "cancel":
		return model.PaymentStatusCanceled, true
	case "deny":
		return model.PaymentStatusFailed, true
	}
	return "", false
}

func bookingTransition(kind model.PaymentKind, target model.PaymentStatus) (from, to model.BookingStatus, ok bool) {
	switch target {
	case model.PaymentStatusPaid:
		switch kind {
		case model.PaymentKindDeposit:
			return model.BookingStatusPending, model.BookingStatusPartiallyPaid, true
		case model.PaymentKindFull:
			return model.BookingStatusPending, model.BookingStatusPaid, true
		case model.PaymentKindBalance:
			return model.BookingStatusPartiallyPaid, model.BookingStatusPaid, true
		}

	case model.PaymentStatusExpired, model.PaymentStatusFailed, model.PaymentStatusCanceled:
		if kind == model.PaymentKindBalance {
			return "", "", false
		}

		switch target {
		case model.PaymentStatusExpired:
			return model.BookingStatusPending, model.BookingStatusExpired, true
		default:
			return model.BookingStatusPending, model.BookingStatusCanceled, true
		}
	}

	return "", "", false
}

func eventBookingTransition(kind model.PaymentKind, target model.PaymentStatus) (from, to model.EventBookingStatus, ok bool) {
	if kind == model.PaymentKindFull && target == model.PaymentStatusPaid {
		return model.EventBookingStatusPending, model.EventBookingStatusPaid, true
	}
	return "", "", false
}
