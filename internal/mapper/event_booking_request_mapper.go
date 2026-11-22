package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToEventBookingRequestResponse(request *model.EventBookingRequest) dto.EventBookingRequestResponse {
	return dto.EventBookingRequestResponse{
		ID:             request.ID,
		User:           ToUserResponse(&request.User),
		Title:          request.Title,
		Description:    request.Description,
		GuestCount:     request.GuestCount,
		Budget:         request.Budget,
		StartAt:        request.StartAt,
		EndAt:          request.EndAt,
		Status:         request.Status,
		ReviewedAt:     request.ReviewedAt,
		RejectReason:   request.RejectReason,
		EventBookingID: request.EventBookingID,
		CreatedAt:      request.CreatedAt,
		UpdatedAt:      request.UpdatedAt,
	}
}

func ToEventBookingRequestModel(req dto.CreateEventRequest, user model.User, villaPackage *model.VillaPackage, start time.Time, end time.Time) *model.EventBookingRequest {
	var villaPackageID *uint
	var villaPackageModel *model.VillaPackage
	if villaPackage != nil {
		villaPackageID = &villaPackage.ID
		villaPackageModel = villaPackage
	}

	return &model.EventBookingRequest{
		UserID:         user.ID,
		User:           user,
		VillaPackageID: villaPackageID,
		VillaPackage:   villaPackageModel,
		Title:          req.Title,
		Description:    req.Description,
		GuestCount:     req.GuestCount,
		Budget:         req.Budget,
		StartAt:        start,
		EndAt:          end,
		Status:         model.RequestPending,
		ReviewedAt:     nil,
		RejectReason:   nil,
		EventBookingID: nil,
	}
}
