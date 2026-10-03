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

type LogoService struct {
	txManager         *repository.TxManager
	repo              *repository.LogoRepository
	specificationRepo *repository.SpecificationRepository
	facilityRepo      *repository.FacilityRepository
}

func NewLogoService(
	txManager *repository.TxManager,
	repo *repository.LogoRepository,
	specificationRepo *repository.SpecificationRepository,
	facilityRepo *repository.FacilityRepository,
) *LogoService {
	return &LogoService{
		txManager:         txManager,
		repo:              repo,
		specificationRepo: specificationRepo,
		facilityRepo:      facilityRepo,
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

func (s *LogoService) GetByID(ctx context.Context, id uint) (dto.LogoResponse, error) {
	logo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.LogoResponse{}, err
	}
	if logo == nil {
		return dto.LogoResponse{}, errors.NotFound(fmt.Sprintf("Logo not found with ID: %d", id))
	}

	return mapper.ToLogoResponse(logo), nil
}

func (s *LogoService) Create(ctx context.Context, req dto.LogoRequest) (dto.LogoResponse, error) {
	var createdLogo *model.Logo

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		logo := mapper.ToLogoModel(req)
		if err := txRepo.Create(ctx, logo); err != nil {
			return err
		}

		createdLogo = logo

		return nil
	}); err != nil {
		return dto.LogoResponse{}, err
	}

	return mapper.ToLogoResponse(createdLogo), nil
}

func (s *LogoService) Update(ctx context.Context, id uint, req dto.LogoRequest) (dto.LogoResponse, error) {
	var updatedLogo *model.Logo

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		logo, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if logo == nil {
			return errors.NotFound(fmt.Sprintf("Logo not found with ID: %d", id))
		}

		mapper.UpdateLogoModel(logo, req)
		if err := txRepo.Update(ctx, logo); err != nil {
			return err
		}

		updatedLogo = logo

		return nil
	}); err != nil {
		return dto.LogoResponse{}, err
	}

	return mapper.ToLogoResponse(updatedLogo), nil
}

func (s *LogoService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txSpecRepo := s.specificationRepo.WithTx(tx)
		txFacRepo := s.facilityRepo.WithTx(tx)

		logo, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if logo == nil {
			return errors.NotFound(fmt.Sprintf("Logo not found with ID: %d", id))
		}

		specExist, err := txSpecRepo.ExistByLogoID(ctx, logo.ID)
		if err != nil {
			return err
		}

		facExist, err := txFacRepo.ExistByLogoID(ctx, logo.ID)
		if err != nil {
			return err
		}

		if facExist || specExist {
			return errors.Conflict("Logo is still used by specifications or facilities")
		}

		return txRepo.Delete(ctx, logo)
	})
}
