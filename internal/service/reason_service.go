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

type ReasonService struct {
	txManager *repository.TxManager
	repo      *repository.ReasonRepository
}

func NewReasonService(
	txManager *repository.TxManager,
	repo *repository.ReasonRepository,
) *ReasonService {
	return &ReasonService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *ReasonService) GetAll(ctx context.Context, query dto.ReasonQuery) ([]dto.ReasonResponse, int64, error) {
	reasons, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	reasonResponses := make([]dto.ReasonResponse, 0, len(reasons))
	for _, reason := range reasons {
		reasonResponses = append(reasonResponses, mapper.ToReasonResponse(&reason))
	}

	return reasonResponses, total, nil
}

func (s *ReasonService) GetByID(ctx context.Context, id uint) (dto.ReasonResponse, error) {
	reason, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.ReasonResponse{}, err
	}
	if reason == nil {
		return dto.ReasonResponse{}, appError.BadRequest(fmt.Sprintf("Reason not found with ID: %d", id))
	}

	return mapper.ToReasonResponse(reason), nil
}

func (s *ReasonService) Create(ctx context.Context, req dto.ReasonRequest) (dto.ReasonResponse, error) {
	var createdReason *model.Reason

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		reason := mapper.ToReasonModel(req)
		if err := txRepo.Create(ctx, reason); err != nil {
			return err
		}

		createdReason = reason

		return nil
	}); err != nil {
		return dto.ReasonResponse{}, err
	}

	return mapper.ToReasonResponse(createdReason), nil
}

func (s *ReasonService) Update(ctx context.Context, req dto.ReasonRequest, id uint) (dto.ReasonResponse, error) {
	var updatedReason *model.Reason

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		reason, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if reason == nil {
			return appError.NotFound(fmt.Sprintf("Reason not found with ID: %d", id))
		}

		mapper.UpdateReasonModel(reason, req)
		if err := txRepo.Update(ctx, reason); err != nil {
			return err
		}

		updatedReason = reason

		return nil
	}); err != nil {
		return dto.ReasonResponse{}, err
	}

	return mapper.ToReasonResponse(updatedReason), nil
}

func (s *ReasonService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		reason, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if reason == nil {
			return appError.NotFound(fmt.Sprintf("Reason not found with ID: %d", id))
		}

		return txRepo.Delete(ctx, reason)
	})
}
