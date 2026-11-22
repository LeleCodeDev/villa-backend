package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToEventBookingResponse(eventBooking *model.EventBooking, payment *model.Payment) dto.EventBookingResponse {
	var villaPackage *dto.VillaPackageResponse
	if eventBooking.VillaPackage != nil {
		villaPackageResponse := ToVillaPackageResponse(eventBooking.VillaPackage, nil)
		villaPackage = &villaPackageResponse
	}

	var paymentResponse *dto.PaymentResponse
	if payment != nil {
		pr := ToPaymentResponse(payment)
		paymentResponse = &pr
	}

	return dto.EventBookingResponse{
		ID:           eventBooking.ID,
		User:         ToUserResponse(&eventBooking.User),
		VillaPackage: villaPackage,
		Title:        eventBooking.Title,
		Status:       eventBooking.Status,
		StartAt:      eventBooking.StartAt,
		EndAt:        eventBooking.EndAt,
		GuestCount:   eventBooking.GuestCount,
		Quantity:     eventBooking.Quantity,
		Price:        eventBooking.Price,
		TotalPrice:   eventBooking.TotalPrice,
		Payment:      paymentResponse,
		CheckInAt:    eventBooking.CheckInAt,
		CheckOutAt:   eventBooking.CheckOutAt,
		ExpiresAt:    eventBooking.ExpiresAt,
		CreatedAt:    eventBooking.CreatedAt,
		UpdatedAt:    eventBooking.UpdatedAt,
	}
}

func ToEventBookingModel(
	startAt time.Time,
	endAt time.Time,
	title string,
	user model.User,
	villaPackage *model.VillaPackage,
	guestCount int,
	quantity int,
	price int64,
	totalPrice int64,
	expireAt time.Time,
) *model.EventBooking {
	var villaPackageID *uint
	if villaPackage != nil {
		villaPackageID = &villaPackage.ID
	}

	return &model.EventBooking{
		UserID:         user.ID,
		User:           user,
		VillaPackageID: villaPackageID,
		VillaPackage:   villaPackage,
		Title:          title,
		Status:         model.EventBookingStatusPending,
		GuestCount:     guestCount,
		StartAt:        startAt,
		EndAt:          endAt,
		Quantity:       quantity,
		Price:          price,
		TotalPrice:     totalPrice,
		CheckInAt:      nil,
		CheckOutAt:     nil,
		ExpiresAt:      &expireAt,
	}
}
