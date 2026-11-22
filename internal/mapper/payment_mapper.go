package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToPaymentResponse(payment *model.Payment) dto.PaymentResponse {
	return dto.PaymentResponse{
		ID:             payment.ID,
		OrderID:        payment.OrderID,
		BookingID:      payment.BookingID,
		EventBookingID: payment.EventBookingID,
		Kind:           payment.Kind,
		Status:         payment.Status,
		Amount:         payment.Amount,
		SnapToken:      payment.SnapToken,
		RedirectURL:    payment.RedirectURL,
		PaidAt:         payment.PaidAt,
		ExpiresAt:      payment.ExpiresAt,
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}
}
