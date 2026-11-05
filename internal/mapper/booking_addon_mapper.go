package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToBookingAddonResponse(bookingAddon *model.BookingAddon) dto.BookingAddonResponse {
	return dto.BookingAddonResponse{
		ID:        bookingAddon.ID,
		BookingID: bookingAddon.BookingID,
		AddonID:   bookingAddon.AddonID,
		Addon:     ToAddonResponse(&bookingAddon.Addon, nil),
		Quantity:  bookingAddon.Quantity,
		CreatedAt: bookingAddon.CreatedAt,
		UpdatedAt: bookingAddon.UpdatedAt,
	}
}

func ToBookingAddonModel(req dto.BookingAddonRequest, booking model.Booking, addon model.Addon) *model.BookingAddon {
	return &model.BookingAddon{
		BookingID: booking.ID,
		Booking:   booking,
		AddonID:   addon.ID,
		Addon:     addon,
		Quantity:  req.Quantity,
	}
}
