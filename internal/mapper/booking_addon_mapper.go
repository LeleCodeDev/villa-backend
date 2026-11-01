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
