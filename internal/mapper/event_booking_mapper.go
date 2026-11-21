package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToEventBookingResponse(eventBooking *model.EventBooking) dto.EventBookingResponse {
	var villaPackage *dto.VillaPackageResponse
	if eventBooking.VillaPackage != nil {
		villaPackageResponse := ToVillaPackageResponse(eventBooking.VillaPackage, nil)
		villaPackage = &villaPackageResponse
	}

	return dto.EventBookingResponse{
		ID:             eventBooking.ID,
		PaymentOrderID: eventBooking.PaymentOrderID,
		User:           ToUserResponse(&eventBooking.User),
		VillaPackage:   villaPackage,
		Status:         eventBooking.Status,
		StartAt:        eventBooking.StartAt,
		EndAt:          eventBooking.EndAt,
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

func ToEventBookingModel(
	startAt time.Time,
	endAt time.Time,
	user model.User,
	villaPackage model.VillaPackage,
	quantity int,
	totalPrice int64,

	expireAt time.Time,
) *model.EventBooking {
	return &model.EventBooking{
		PaymentOrderID: nil,
		UserID:         user.ID,
		User:           user,
		VillaPackageID: &villaPackage.ID,
		VillaPackage:   &villaPackage,
		Status:         model.BookingStatusPending,
		StartAt:        startAt,
		EndAt:          endAt,
		Quantity:       quantity,
		TotalPrice:     totalPrice,
		SnapToken:      nil,
		RedirectURL:    nil,
		PaidAt:         nil,
		CheckInAt:      nil,
		CheckOutAt:     nil,
		ExpiresAt:      &expireAt,
	}
}
