package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/repository"
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
