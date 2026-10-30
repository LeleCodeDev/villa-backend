package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToBookingResponse(booking *model.Booking, addons []model.Addon) dto.BookingResponse {
	var addonResponses []dto.AddonResponse
	for i := range addons {
		addonResponses = append(addonResponses, ToAddonResponse(&addons[i], nil))
	}

	return dto.BookingResponse{
		ID:                    booking.ID,
		User:                  ToUserResponse(&booking.User),
		GuestCount:            booking.GuestCount,
		RentalOption:          ToRentalOptionResponse(&booking.RentalOption, nil),
		Reason:                ToReasonResponse(&booking.Reason),
		Status:                booking.Status,
		CheckInDate:           booking.CheckInDate,
		CheckOutDate:          booking.CheckOutDate,
		BookingPrice:          booking.BookingPrice,
		TotalPrice:            booking.TotalPrice,
		CheckInAt:             booking.CheckInAt,
		CheckOutAt:            &booking.CheckOutDate,
		ConfirmationExpiresAt: booking.ConfirmationExpiresAt,
		Addons:                addonResponses,
		CreatedAt:             booking.CreatedAt,
		UpdatedAt:             booking.UpdatedAt,
	}
}
