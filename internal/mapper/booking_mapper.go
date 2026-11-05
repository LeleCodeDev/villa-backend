package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToBookingResponse(booking *model.Booking, bookingAddons []model.BookingAddon) dto.BookingResponse {
	var BookingAddonResponses []dto.BookingAddonResponse
	for i := range bookingAddons {
		BookingAddonResponses = append(BookingAddonResponses, ToBookingAddonResponse(&bookingAddons[i]))
	}

	var rentalOption *dto.RentalOptionResponse
	if booking.RentalOption != nil {
		rentalOptionResponse := ToRentalOptionResponse(booking.RentalOption, nil)
		rentalOption = &rentalOptionResponse
	}

	return dto.BookingResponse{
		ID:             booking.ID,
		PaymentOrderID: booking.PaymentOrderID,
		User:           ToUserResponse(&booking.User),
		GuestCount:     booking.GuestCount,
		RentalOption:   rentalOption,
		Reason:         ToReasonResponse(&booking.Reason),
		Status:         booking.Status,
		CheckInDate:    booking.CheckInDate,
		CheckOutDate:   booking.CheckOutDate,
		BookingPrice:   booking.BookingPrice,
		TotalPrice:     booking.TotalPrice,
		SnapToken:      booking.SnapToken,
		RedirectURL:    booking.RedirectURL,
		PaidAt:         booking.PaidAt,
		CheckInAt:      booking.CheckInAt,
		CheckOutAt:     booking.CheckOutAt,
		ExpiresAt:      booking.ExpiresAt,
		BookingAddons:  BookingAddonResponses,
		CreatedAt:      booking.CreatedAt,
		UpdatedAt:      booking.UpdatedAt,
	}
}

func ToBookingModel(
	req dto.BookingCreateRequest,
	user model.User,
	rentalOption model.RentalOption,
	reason model.Reason,
	checkInDate time.Time,
	checkOutDate time.Time,
	bookingPrice int64,
	totalPrice int64,
	expireAt time.Time,
) *model.Booking {
	return &model.Booking{
		UserID:         user.ID,
		PaymentOrderID: nil,
		User:           user,
		GuestCount:     req.GuestCount,
		RentalOptionID: &rentalOption.ID,
		RentalOption:   &rentalOption,
		ReasonID:       reason.ID,
		Reason:         reason,
		Status:         model.BookingStatusPending,
		CheckInDate:    checkInDate,
		CheckOutDate:   checkOutDate,
		BookingPrice:   bookingPrice,
		TotalPrice:     totalPrice,
		SnapToken:      nil,
		RedirectURL:    nil,
		PaidAt:         nil,
		CheckInAt:      nil,
		CheckOutAt:     nil,
		ExpiresAt:      &expireAt,
	}
}
