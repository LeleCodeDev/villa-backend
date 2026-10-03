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

type FacilityService struct {
	txManager *repository.TxManager
	repo      *repository.FacilityRepository
	logoRepo  *repository.LogoRepository
}

func NewFacilityService(
	txManager *repository.TxManager,
	repo *repository.FacilityRepository,
	logoRepo *repository.LogoRepository,
) *FacilityService {
	return &FacilityService{
		txManager: txManager,
		repo:      repo,
		logoRepo:  logoRepo,
	}
}

func (s *FacilityService) GetAll(ctx context.Context, query dto.FacilityQuery) ([]dto.FacilityResponse, int64, error) {
	facilities, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.FacilityResponse, 0, len(facilities))
	for _, facility := range facilities {
		responses = append(responses, mapper.ToFacilityResponse(&facility))
	}

	return responses, total, nil
}

func (s *FacilityService) GetByID(ctx context.Context, id uint) (dto.FacilityResponse, error) {
	facility, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.FacilityResponse{}, err
	}
	if facility == nil {
		return dto.FacilityResponse{}, errors.NotFound(fmt.Sprintf("Facility not found with ID: %d", id))
	}

	return mapper.ToFacilityResponse(facility), nil
}

func (s *FacilityService) Create(ctx context.Context, req dto.FacilityRequest) (dto.FacilityResponse, error) {
	var createdFacility *model.Facility

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txLogoRepo := s.logoRepo.WithTx(tx)

		logo, err := txLogoRepo.GetByID(ctx, req.LogoID)
		if err != nil {
			return err
		}
		if logo == nil {
			return errors.NotFound(fmt.Sprintf("Logo not found with ID : %d", req.LogoID))
		}

		var sortOrder int
		if req.SortOrder == nil {
			maxOrder, err := txRepo.GetMaxSortOrder(ctx)
			if err != nil {
				return err
			}

			sortOrder = maxOrder + 1
		} else {
			sortOrder = *req.SortOrder
			if err := txRepo.ShiftSortOrder(ctx, sortOrder); err != nil {
				return err
			}
		}

		facility := mapper.ToFacilityModel(req, sortOrder, logo)
		if err := txRepo.Create(ctx, facility); err != nil {
			return err
		}

		createdFacility = facility

		return nil
	}); err != nil {
		return dto.FacilityResponse{}, err
	}

	return mapper.ToFacilityResponse(createdFacility), nil
}

func (s *FacilityService) Update(ctx context.Context, id uint, req dto.FacilityRequest) (dto.FacilityResponse, error) {
	var updatedFacility *model.Facility

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txLogoRepo := s.logoRepo.WithTx(tx)

		facility, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if facility == nil {
			return errors.NotFound(fmt.Sprintf("Facility not found with ID: %d", id))
		}

		logo := &facility.Logo
		if facility.LogoID != req.LogoID {
			logo, err = txLogoRepo.GetByID(ctx, req.LogoID)
			if err != nil {
				return err
			}
			if logo == nil {
				return errors.NotFound(fmt.Sprintf("Logo not found with ID : %d", req.LogoID))
			}
		}

		oldOrder := facility.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if *newOrder != oldOrder {
			if *newOrder < oldOrder {
				if err := txRepo.UpdateSortOrderRange(ctx, *newOrder, oldOrder-1, +1); err != nil {
					return err
				}
			} else {
				if err := txRepo.UpdateSortOrderRange(ctx, oldOrder+1, *newOrder, -1); err != nil {
					return err
				}
			}
		}

		mapper.UpdateFacilityModel(facility, req, logo, *newOrder)
		if err := txRepo.Update(ctx, facility); err != nil {
			return err
		}

		updatedFacility = facility

		return nil
	}); err != nil {
		return dto.FacilityResponse{}, err
	}

	return mapper.ToFacilityResponse(updatedFacility), nil
}

func (s *FacilityService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		facility, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if facility == nil {
			return errors.NotFound(fmt.Sprintf("Facility not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, facility.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, facility)
	})
}
