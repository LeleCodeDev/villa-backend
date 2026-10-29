package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	BookingQuery struct {
		PaginationQuery
		SortQuery
		Status    model.BookingStatus `json:"status" form:"status" binding:"omitempty,oneof=pending paid canceled completed expired refunded not_show"`
		StartDate string              `json:"start_date" form:"start_date" binding:"datetime=2006-01-02"`
		EndDate   string              `json:"end_date" form:"end_date" binding:"datetime=2006-01-02"`
		Unpage    bool                `form:"unpage"`
	}

	BookingCreateRequest struct {
		GuestCount     int    `json:"guest_count" form:"guest_count" binding:"required,gt=0"`
		RentalOptionID uint   `json:"rental_option_id" form:"rental_option_id" binding:"required,gt=0"`
		ReasonID       uint   `json:"reason_id" form:"reason_id" binding:"required,gt=0"`
		CheckInDate    string `json:"check_in_date" form:"check_in_date" binding:"required,datetime=2006-01-02"`
		CheckOutDate   string `json:"check_out_date" form:"check_out_date" binding:"required,datetime=2006-01-02"`
	}

	BookingResponse struct {
		ID                    int                  `json:"id"`
		User                  UserResponse         `json:"user"`
		GuestCount            int                  `json:"guest_count"`
		RentalOption          RentalOptionResponse `json:"rental_option"`
		Reason                ReasonResponse       `json:"reason"`
		Status                model.BookingStatus  `json:"status"`
		CheckInDate           time.Time            `json:"check_in_date"`
		CheckOutDate          time.Time            `json:"check_out_date"`
		TotalPrice            float64              `json:"total_price"`
		CheckInAt             *time.Time           `json:"check_in_at"`
		CheckOutAt            *time.Time           `json:"check_out_at"`
		ConfirmationExpiresAt *time.Time           `json:"confirmation_expires_at,omitempty"`
		Addons                []AddonResponse      `json:"addons"`
		CreatedAt             time.Time            `json:"created_at"`
		UpdatedAt             time.Time            `json:"updated_at"`
	}
)

func (bq *BookingQuery) SetDefault() {
	bq.setDefaultPagination()
	bq.SetDefaultSort(SortAsc)
}
