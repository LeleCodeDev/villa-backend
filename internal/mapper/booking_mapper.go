package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToBookingResponse(booking *model.Booking, bookingAddons []model.BookingAddon) dto.BookingResponse {
	var BookingAddonResponses []dto.BookingAddonResponse
	for i := range bookingAddons {
		BookingAddonResponses = append(BookingAddonResponses, ToBookingAddonResponse(&bookingAddons[i]))
	}

	return dto.BookingResponse{
		ID:            booking.ID,
		User:          ToUserResponse(&booking.User),
		GuestCount:    booking.GuestCount,
		RentalOption:  ToRentalOptionResponse(&booking.RentalOption, nil),
		Reason:        ToReasonResponse(&booking.Reason),
		Status:        booking.Status,
		CheckInDate:   booking.CheckInDate,
		CheckOutDate:  booking.CheckOutDate,
		BookingPrice:  booking.BookingPrice,
		TotalPrice:    booking.TotalPrice,
		CheckInAt:     booking.CheckInAt,
		CheckOutAt:    &booking.CheckOutDate,
		ExpiresAt:     booking.ExpiresAt,
		BookingAddons: BookingAddonResponses,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
}
