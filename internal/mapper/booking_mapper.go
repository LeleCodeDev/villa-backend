package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToBookingResponse(booking *model.Booking, bookingAddons []model.BookingAddon, payments []model.Payment) dto.BookingResponse {
	var bookingAddonResponses []dto.BookingAddonResponse
	for i := range bookingAddons {
		bookingAddonResponses = append(bookingAddonResponses, ToBookingAddonResponse(&bookingAddons[i]))
	}

	var paymentResponses []dto.PaymentResponse
	for i := range payments {
		paymentResponses = append(paymentResponses, ToPaymentResponse(&payments[i]))
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
		VisitorCount:  booking.VisitorCount,
		RentalOption:  rentalOption,
		Reason:        ToReasonResponse(&booking.Reason),
		Status:        booking.Status,
		PaymentType:   booking.PaymentType,
		CheckInDate:   booking.CheckInDate,
		CheckOutDate:  booking.CheckOutDate,
		BookingPrice:  booking.BookingPrice,
		TotalPrice:    booking.TotalPrice,
		DepositAmount: booking.DepositAmount,
		BalanceDueAt:  booking.BalanceDueAt,
		CheckInAt:     booking.CheckInAt,
		CheckOutAt:    booking.CheckOutAt,
		ExpiresAt:     booking.ExpiresAt,
		BookingAddons: bookingAddonResponses,
		Payments:      paymentResponses,
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
	bookingPrice int64,
	totalPrice int64,
	depositAmount *int64,
	balanceDueAt *time.Time,
	expireAt time.Time,
) *model.Booking {
	return &model.Booking{
		UserID:         user.ID,
		User:           user,
		GuestCount:     req.GuestCount,
		VisitorCount:   req.VisitorCount,
		RentalOptionID: &rentalOption.ID,
		RentalOption:   &rentalOption,
		ReasonID:       reason.ID,
		Reason:         reason,
		Status:         model.BookingStatusPending,
		PaymentType:    req.PaymentType,
		CheckInDate:    checkInDate,
		CheckOutDate:   checkOutDate,
		BookingPrice:   bookingPrice,
		TotalPrice:     totalPrice,
		DepositAmount:  depositAmount,
		BalanceDueAt:   balanceDueAt,
		CheckInAt:      nil,
		CheckOutAt:     nil,
		ExpiresAt:      &expireAt,
	}
}
