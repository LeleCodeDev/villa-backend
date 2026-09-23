package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type VillaPackageListService struct {
	txManager *repository.TxManager
	repo      *repository.VillaPackageListRepository
}

func NewVillaPackageListService(
	txManager *repository.TxManager,
	repo *repository.VillaPackageListRepository,
) *VillaPackageListService {
	return &VillaPackageListService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *VillaPackageListService) GetAll(ctx context.Context) ([]dto.VillaPackageListResponse, error) {
	villaPackageLists, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VillaPackageListResponse, 0, len(villaPackageLists))
	for _, villaPackageList := range villaPackageLists {
		responses = append(responses, mapper.ToVillaPackageListResponse(&villaPackageList))
	}

	return responses, nil
}
