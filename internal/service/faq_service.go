package service

import (
	"context"
	"fmt"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"gorm.io/gorm"
)

type FaqService struct {
	txManager *repository.TxManager
	repo      *repository.FaqRepository
}

func NewFaqService(
	txManager *repository.TxManager,
	repo *repository.FaqRepository,
) *FaqService {
	return &FaqService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *FaqService) GetAll(ctx context.Context) ([]dto.FaqResponse, error) {
	faqs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.FaqResponse, 0, len(faqs))
	for _, faq := range faqs {
		responses = append(responses, mapper.ToFaqResponse(&faq))
	}

	return responses, nil
}

func (s *FaqService) GetByID(ctx context.Context, id uint) (dto.FaqResponse, error) {
	faq, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.FaqResponse{}, err
	}
	if faq == nil {
		return dto.FaqResponse{}, errors.NotFound(fmt.Sprintf("Faq not found with ID: %d", id))
	}

	return mapper.ToFaqResponse(faq), nil
}

func (s *FaqService) Create(ctx context.Context, req dto.FaqRequest) (dto.FaqResponse, error) {
	var createdFaq *model.Faq

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		faq := mapper.ToFaqModel(req)
		if err := txRepo.Create(ctx, faq); err != nil {
			return err
		}

		createdFaq = faq

		return nil
	}); err != nil {
		return dto.FaqResponse{}, err
	}

	return mapper.ToFaqResponse(createdFaq), nil
}

func (s *FaqService) Update(ctx context.Context, id uint, req dto.FaqRequest) (dto.FaqResponse, error) {
	var updatedFaq *model.Faq

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		faq, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if faq == nil {
			return errors.NotFound(fmt.Sprintf("Faq not found with ID: %d", id))
		}

		mapper.UpdateFaqModel(faq, req)
		if err := txRepo.Update(ctx, faq); err != nil {
			return err
		}

		updatedFaq = faq

		return nil
	}); err != nil {
		return dto.FaqResponse{}, err
	}

	return mapper.ToFaqResponse(updatedFaq), nil
}
