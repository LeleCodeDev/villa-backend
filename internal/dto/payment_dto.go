package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	MidtransCallbackRequest struct {
		OrderID           string `json:"order_id"`
		StatusCode        string `json:"status_code"`
		GrossAmount       string `json:"gross_amount"`
		SignatureKey      string `json:"signature_key"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
		PaymentType       string `json:"payment_type"`
	}

	PaymentResponse struct {
		ID             uint                `json:"id"`
		OrderID        string              `json:"order_id"`
		BookingID      *uint               `json:"booking_id,omitempty"`
		EventBookingID *uint               `json:"event_booking_id,omitempty"`
		Kind           model.PaymentKind   `json:"kind"`
		Status         model.PaymentStatus `json:"status"`
		Amount         int64               `json:"amount"`
		SnapToken      *string             `json:"snap_token"`
		RedirectURL    *string             `json:"redirect_url"`
		PaidAt         *time.Time          `json:"paid_at"`
		ExpiresAt      *time.Time          `json:"expires_at,omitempty"`
		CreatedAt      time.Time           `json:"created_at"`
		UpdatedAt      time.Time           `json:"updated_at"`
	}
)
