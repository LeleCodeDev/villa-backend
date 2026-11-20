package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type EventBookingRequestService struct {
	txManager *repository.TxManager
	repo      *repository.EventBookingRequestRepository
}

func NewEventBookingRequestService(
	txManager *repository.TxManager,
	repo *repository.EventBookingRequestRepository,
) *EventBookingRequestService {
	return &EventBookingRequestService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *EventBookingRequestService) GetAll(ctx context.Context, query dto.EventBookingRequestQuery, userID *uint) ([]dto.EventBookingRequestResponse, int64, error) {
	requests, total, err := s.repo.GetAll(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.EventBookingRequestResponse, 0, len(requests))
	for _, request := range requests {
		responses = append(responses, mapper.ToEventBookingRequestResponse(&request))
	}

	return responses, total, nil
}

func (s *EventBookingRequestService) GetByID(ctx context.Context, id uint, currentUser model.User) (dto.EventBookingRequestResponse, error) {
	request, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.EventBookingRequestResponse{}, err
	}
	if request == nil {
		return dto.EventBookingRequestResponse{}, appError.NotFound(fmt.Sprintf("Event request not found with ID: %d", id))
	}

	if currentUser.Role != model.RoleAdmin && request.UserID != currentUser.ID {
		return dto.EventBookingRequestResponse{}, appError.Unauthorized("You are not authorized to access this event request!")
	}

	return mapper.ToEventBookingRequestResponse(request), nil
}

func (s *EventBookingRequestService) Create(ctx context.Context, req dto.CreateEventRequest, currentUser model.User) (dto.EventBookingRequestResponse, error) {
	startAt, endAt, err := parseEventBookingDates(req.StartAt, req.EndAt)
	if err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	var createdRequest *model.EventBookingRequest

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		request := mapper.ToEventBookingRequestModel(req, currentUser, startAt, endAt)
		if err := txRepo.Create(ctx, request); err != nil {
			return err
		}

		createdRequest = request

		return nil
	}); err != nil {
		return dto.EventBookingRequestResponse{}, err
	}

	return mapper.ToEventBookingRequestResponse(createdRequest), nil
}

// func (s *EventBookingRequestService) Approve(ctx context.Context, req dto.ApproveEventRequest) (dto.EventBookingRequestResponse, error) {
// }
