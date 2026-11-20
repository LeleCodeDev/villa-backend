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

func ToEventBookingRequestModel(req dto.CreateEventRequest, user model.User, start time.Time, end time.Time) *model.EventBookingRequest {
	return &model.EventBookingRequest{
		UserID:         user.ID,
		User:           user,
		Title:          req.Title,
		Description:    req.Description,
		StartAt:        start,
		EndAt:          end,
		Status:         model.RequestPending,
		ReviewedAt:     nil,
		RejectReason:   nil,
		EventBookingID: nil,
	}
}
