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

type SpecificationService struct {
	txManager *repository.TxManager
	repo      *repository.SpecificationRepository
	logoRepo  *repository.LogoRepository
}

func NewSpecificationService(
	txManager *repository.TxManager,
	repo *repository.SpecificationRepository,
	logoRepo *repository.LogoRepository,
) *SpecificationService {
	return &SpecificationService{
		txManager: txManager,
		repo:      repo,
		logoRepo:  logoRepo,
	}
}

func (s *SpecificationService) GetAll(ctx context.Context, query dto.SpecificationQuery) ([]dto.SpecificationResponse, int64, error) {
	specifications, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.SpecificationResponse, 0, len(specifications))
	for _, specification := range specifications {
		responses = append(responses, mapper.ToSpecificationResponse(&specification))
	}

	return responses, total, nil
}

func (s *SpecificationService) GetByID(ctx context.Context, id uint) (dto.SpecificationResponse, error) {
	specification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.SpecificationResponse{}, err
	}
	if specification == nil {
		return dto.SpecificationResponse{}, errors.NotFound(fmt.Sprintf("Specification not found with ID: %d", id))
	}

	return mapper.ToSpecificationResponse(specification), nil
}

func (s *SpecificationService) Create(ctx context.Context, req dto.SpecificationRequest) (dto.SpecificationResponse, error) {
	var createdSpecification *model.Specification

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

		specification := mapper.ToSpecificationModel(req, sortOrder, logo)
		if err := txRepo.Create(ctx, specification); err != nil {
			return err
		}

		createdSpecification = specification

		return nil
	}); err != nil {
		return dto.SpecificationResponse{}, err
	}

	return mapper.ToSpecificationResponse(createdSpecification), nil
}

func (s *SpecificationService) Update(ctx context.Context, id uint, req dto.SpecificationRequest) (dto.SpecificationResponse, error) {
	var updatedSpecification *model.Specification

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txLogoRepo := s.logoRepo.WithTx(tx)

		specification, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if specification == nil {
			return errors.NotFound(fmt.Sprintf("Specification not found with ID: %d", id))
		}

		logo := &specification.Logo
		if specification.LogoID != req.LogoID {
			logo, err = txLogoRepo.GetByID(ctx, req.LogoID)
			if err != nil {
				return err
			}
			if logo == nil {
				return errors.NotFound(fmt.Sprintf("Logo not found with ID : %d", req.LogoID))
			}
		}

		oldOrder := specification.SortOrder
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

		mapper.UpdateSpecificationModel(specification, req, logo, *newOrder)
		if err := txRepo.Update(ctx, specification); err != nil {
			return err
		}

		updatedSpecification = specification

		return nil
	}); err != nil {
		return dto.SpecificationResponse{}, err
	}

	return mapper.ToSpecificationResponse(updatedSpecification), nil
}

func (s *SpecificationService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		specification, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if specification == nil {
			return errors.NotFound(fmt.Sprintf("Specification not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, specification.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, specification)
	})
}
