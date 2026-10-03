package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type LogoService struct {
	txManager *repository.TxManager
	repo      *repository.LogoRepository
}

func NewLogoService(
	txManager *repository.TxManager,
	repo *repository.LogoRepository,
) *LogoService {
	return &LogoService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *LogoService) GetAll(ctx context.Context, query dto.LogoQuery) ([]dto.LogoResponse, int64, error) {
	logos, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.LogoResponse, 0, len(logos))
	for _, logo := range logos {
		responses = append(responses, mapper.ToLogoResponse(&logo))
	}

	return responses, total, nil
}
