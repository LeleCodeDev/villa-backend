package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToEventBookingResponse(eventBooking *model.EventBooking) dto.EventBookingResponse {
	return dto.EventBookingResponse{
		ID:             eventBooking.ID,
		PaymentOrderID: eventBooking.PaymentOrderID,
		User:           ToUserResponse(&eventBooking.User),
		Status:         eventBooking.Status,
		Type:           eventBooking.Type,
		StartAt:        eventBooking.StartAt,
		EndAt:          eventBooking.EndAt,
		Price:          eventBooking.Price,
		Quantity:       eventBooking.Quantity,
		TotalPrice:     eventBooking.TotalPrice,
		SnapToken:      eventBooking.SnapToken,
		RedirectURL:    eventBooking.RedirectURL,
		PaidAt:         eventBooking.PaidAt,
		CheckInAt:      eventBooking.CheckInAt,
		CheckOutAt:     eventBooking.CheckOutAt,
		ExpiresAt:      eventBooking.ExpiresAt,
		CreatedAt:      eventBooking.CreatedAt,
		UpdatedAt:      eventBooking.UpdatedAt,
	}
}
