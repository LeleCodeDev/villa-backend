package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type TestimonyService struct {
	txManager *repository.TxManager
	repo      *repository.TestimonyRepository
}

func NewTestimonyService(
	txManager *repository.TxManager,
	repo *repository.TestimonyRepository,
) *TestimonyService {
	return &TestimonyService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *TestimonyService) GetAll(ctx context.Context) ([]dto.TestimonyResponse, error) {
	testimonies, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TestimonyResponse, 0, len(testimonies))
	for _, testimony := range testimonies {
		responses = append(responses, mapper.ToTestimonyResponse(&testimony))
	}

	return responses, nil
}

func (s *TestimonyService) GetByID(ctx context.Context, id uint) (dto.TestimonyResponse, error) {
	testimony, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.TestimonyResponse{}, err
	}
	if testimony == nil {
		return dto.TestimonyResponse{}, err
	}

	return mapper.ToTestimonyResponse(testimony), nil
}
