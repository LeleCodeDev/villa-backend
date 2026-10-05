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

type MessageService struct {
	txManager *repository.TxManager
	repo      *repository.MessageRepository
}

func NewMessageService(
	txManager *repository.TxManager,
	repo *repository.MessageRepository,
) *MessageService {
	return &MessageService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *MessageService) GetAll(ctx context.Context, query dto.MessageQuery) ([]dto.MessageResponse, int64, error) {
	messages, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.MessageResponse, 0, len(messages))
	for _, message := range messages {
		responses = append(responses, mapper.ToMessageResponse(&message))
	}

	return responses, total, nil
}

func (s *MessageService) GetByID(ctx context.Context, id uint) (dto.MessageResponse, error) {
	message, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.MessageResponse{}, err
	}
	if message == nil {
		return dto.MessageResponse{}, appError.NotFound(fmt.Sprintf("Message not found with ID: %d", id))
	}

	return mapper.ToMessageResponse(message), nil
}

func (s *MessageService) Create(ctx context.Context, req dto.MessageRequest) (dto.MessageResponse, error) {
	var createdMessage *model.Message

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		message := mapper.ToMessageModel(req)
		if err := txRepo.Create(ctx, message); err != nil {
			return err
		}

		createdMessage = message

		return nil
	}); err != nil {
		return dto.MessageResponse{}, err
	}

	return mapper.ToMessageResponse(createdMessage), nil
}
