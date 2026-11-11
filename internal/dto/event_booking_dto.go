package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	EventBookingQuery struct {
		PaginationQuery
		SortQuery
		Status  model.BookingStatus    `form:"status" binding:"omitempty,oneof=pending paid canceled completed expired refunded not_show"`
		Type    model.EventBookingType `form:"type" binding:"omitempty,oneof=per_hour per_day one_time"`
		StartAt string                 `form:"start_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		EndAt   string                 `form:"end_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		Unpage  bool                   `form:"unpage"`
	}

	EventBookingCreateRequest struct {
		GuestCount         int                   `json:"guest_count" form:"guest_count" binding:"required,gt=0"`
		RentalOptionID     uint                  `json:"rental_option_id" form:"rental_option_id" binding:"required,gt=0"`
		ReasonID           uint                  `json:"reason_id" form:"reason_id" binding:"required,gt=0"`
		CheckInDate        string                `json:"check_in_date" form:"check_in_date" binding:"required,datetime=2006-01-02"`
		CheckOutDate       string                `json:"check_out_date" form:"check_out_date" binding:"required,datetime=2006-01-02"`
		EventBookingAddons []BookingAddonRequest `json:"booking_addons" form:"booking_addons" binding:"omitempty,dive"`
	}

	EventBookingResponse struct {
		ID             uint                   `json:"id"`
		PaymentOrderID *string                `json:"payment_order_id"`
		User           UserResponse           `json:"user"`
		Status         model.BookingStatus    `json:"status"`
		Type           model.EventBookingType `json:"type"`
		StartAt        time.Time              `json:"start_at"`
		EndAt          time.Time              `json:"end_at"`
		Price          int64                  `json:"price"`
		Quantity       int                    `json:"quantity"`
		TotalPrice     int64                  `json:"total_price"`
		SnapToken      *string                `json:"snap_token"`
		RedirectURL    *string                `json:"redirect_url,omitempty"`
		PaidAt         *time.Time             `json:"paid_at"`
		CheckInAt      *time.Time             `json:"check_in_at"`
		CheckOutAt     *time.Time             `json:"check_out_at"`
		ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
		CreatedAt      time.Time              `json:"created_at"`
		UpdatedAt      time.Time              `json:"updated_at"`
	}

	EventBookingAvailableRequest struct {
		CheckInDate  string `json:"check_in_date" form:"check_in_date" binding:"required,datetime=2006-01-02"`
		CheckOutDate string `json:"check_out_date" form:"check_out_date" binding:"required,datetime=2006-01-02"`
		GuestCount   int    `json:"guest_count" form:"guest_count" binding:"required,gt=0"`
	}

	EventBookingConflict struct {
		CheckInDate  time.Time `json:"check_in_date"`
		CheckOutDate time.Time `json:"check_out_date"`
	}
)

func (bq *EventBookingQuery) SetDefault() {
	bq.setDefaultPagination()
	bq.SetDefaultSort(SortAsc)
}
