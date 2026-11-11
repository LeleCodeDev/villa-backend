package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	EventBookingQuery struct {
		PaginationQuery
		SortQuery
		Status  model.BookingStatus `form:"status" binding:"omitempty,oneof=pending paid canceled completed expired refunded not_show"`
		Type    model.ChargeType    `form:"type" binding:"omitempty,oneof=per_hour per_day one_time"`
		StartAt string              `form:"start_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		EndAt   string              `form:"end_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		Unpage  bool                `form:"unpage"`
	}

	EventBookingCreateRequest struct {
		StartAt        string `json:"start_at" form:"start_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		EndAt          string `json:"end_at" form:"end_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		VillaPackageID uint   `json:"villa_package_id" form:"villa_package_id" binding:"required,gt=0"`
	}

	EventBookingResponse struct {
		ID             uint                  `json:"id"`
		PaymentOrderID *string               `json:"payment_order_id"`
		User           UserResponse          `json:"user"`
		VillaPackage   *VillaPackageResponse `json:"villa_package"`
		Status         model.BookingStatus   `json:"status"`
		StartAt        time.Time             `json:"start_at"`
		EndAt          time.Time             `json:"end_at"`
		Price          int64                 `json:"price"`
		Quantity       int                   `json:"quantity"`
		TotalPrice     int64                 `json:"total_price"`
		SnapToken      *string               `json:"snap_token"`
		RedirectURL    *string               `json:"redirect_url,omitempty"`
		PaidAt         *time.Time            `json:"paid_at"`
		CheckInAt      *time.Time            `json:"check_in_at"`
		CheckOutAt     *time.Time            `json:"check_out_at"`
		ExpiresAt      *time.Time            `json:"expires_at,omitempty"`
		CreatedAt      time.Time             `json:"created_at"`
		UpdatedAt      time.Time             `json:"updated_at"`
	}

	EventBookingAvailableRequest struct {
		CheckInDate  string `json:"check_in_date" form:"check_in_date" binding:"required,datetime=2006-01-02"`
		CheckOutDate string `json:"check_out_date" form:"check_out_date" binding:"required,datetime=2006-01-02"`
		GuestCount   int    `json:"guest_count" form:"guest_count" binding:"required,gt=0"`
	}

	EventBookingConflict struct {
		StartAt time.Time `json:"start_at"`
		EndAt   time.Time `json:"end_at"`
	}
)

func (bq *EventBookingQuery) SetDefault() {
	bq.setDefaultPagination()
	bq.SetDefaultSort(SortAsc)
}
