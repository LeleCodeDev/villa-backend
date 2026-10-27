package dto

import "time"

type (
	BookingCreateRequest struct {
		GuestCount     int       `json:"guest_count" form:"guest_count" binding:"required,gt=0"`
		RentalOptionID uint      `json:"rental_option_id" form:"rental_option_id" binding:"required,gt=0"`
		ReasonID       uint      `json:"reason_id" form:"reason_id" binding:"required,gt=0"`
		CheckInDate    time.Time `gorm:"not null;type:date"`
		CheckOutDate   time.Time `gorm:"not null;type:date"`

		// ReturnDate *types.DateOnly `json:"returnDate" binding:"required" time_format:"2006-01-02"`
	}
)
