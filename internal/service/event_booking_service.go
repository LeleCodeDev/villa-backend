package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/payment"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type EventBookingService struct {
	txManager *repository.TxManager
	repo      *repository.EventBookingRepository
	gateaway  *payment.MidtransGateway
}

func NewEventBookingService(
	txManager *repository.TxManager,
	repo *repository.EventBookingRepository,
	gateaway *payment.MidtransGateway,
) *EventBookingService {
	return &EventBookingService{
		txManager: txManager,
		repo:      repo,
		gateaway:  gateaway,
	}
}

func (s *EventBookingService) GetAll(ctx context.Context, query dto.EventBookingQuery, userID *uint) ([]dto.EventBookingResponse, int64, error) {
	eventBokings, total, err := s.repo.GetAll(ctx, query, userID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.EventBookingResponse, 0, len(eventBokings))
	for _, eventBoking := range eventBokings {
		responses = append(responses, mapper.ToEventBookingResponse(&eventBoking))
	}

	return responses, total, nil
}
