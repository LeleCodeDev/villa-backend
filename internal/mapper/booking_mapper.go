package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/config"
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
		ID:            booking.ID,
		User:          ToUserResponse(&booking.User),
		GuestCount:    booking.GuestCount,
		RentalOption:  rentalOption,
		Reason:        ToReasonResponse(&booking.Reason),
		Status:        booking.Status,
		CheckInDate:   booking.CheckInDate,
		CheckOutDate:  booking.CheckOutDate,
		BookingPrice:  booking.BookingPrice,
		TotalPrice:    booking.TotalPrice,
		CheckInAt:     booking.CheckInAt,
		CheckOutAt:    booking.CheckOutAt,
		ExpiresAt:     booking.ExpiresAt,
		BookingAddons: BookingAddonResponses,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}
}

func ToBookingModel(
	req dto.BookingCreateRequest,
	user model.User,
	rentalOption model.RentalOption,
	reason model.Reason,
	checkInDate time.Time,
	checkOutDate time.Time,
	bookingPrice float64,
	totalPrice float64,
) *model.Booking {
	expireAt := time.Now().In(config.Env.Location).Add(time.Hour)

	return &model.Booking{
		UserID:         user.ID,
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
		CheckInAt:      nil,
		CheckOutAt:     nil,
		ExpiresAt:      &expireAt,
	}
}
