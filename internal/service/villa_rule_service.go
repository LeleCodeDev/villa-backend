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

type VillaRuleService struct {
	txManager *repository.TxManager
	repo      *repository.VillaRuleRepository
}

func NewVillaRuleService(
	txManager *repository.TxManager,
	repo *repository.VillaRuleRepository,
) *VillaRuleService {
	return &VillaRuleService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *VillaRuleService) GetAll(ctx context.Context, query dto.VillaRuleQuery) ([]dto.VillaRuleResponse, int64, error) {
	villaRules, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.VillaRuleResponse, 0, len(villaRules))
	for _, villaRule := range villaRules {
		responses = append(responses, mapper.ToVillaRuleResponse(&villaRule))
	}

	return responses, total, nil
}

func (s *VillaRuleService) GetByID(ctx context.Context, id uint) (dto.VillaRuleResponse, error) {
	villaRule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.VillaRuleResponse{}, err
	}
	if villaRule == nil {
		return dto.VillaRuleResponse{}, appError.NotFound(fmt.Sprintf("Villa rule not found with ID: %d", id))
	}

	return mapper.ToVillaRuleResponse(villaRule), nil
}

func (s *VillaRuleService) Create(ctx context.Context, req dto.VillaRuleRequest) (dto.VillaRuleResponse, error) {
	var createdVillaRule *model.VillaRule

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

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

		villaRule := mapper.ToVillaRuleModel(req, sortOrder)
		if err := txRepo.Create(ctx, villaRule); err != nil {
			return err
		}

		createdVillaRule = villaRule

		return nil
	}); err != nil {
		return dto.VillaRuleResponse{}, err
	}

	return mapper.ToVillaRuleResponse(createdVillaRule), nil
}

func (s *VillaRuleService) Update(ctx context.Context, id uint, req dto.VillaRuleRequest) (dto.VillaRuleResponse, error) {
	var updatedVillaRule *model.VillaRule

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		villaRule, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaRule == nil {
			return appError.NotFound(fmt.Sprintf("Villa rule not found with ID: %d", id))
		}

		oldOrder := villaRule.SortOrder
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

		mapper.UpdateVillaRuleModel(villaRule, req, *newOrder)
		if err := txRepo.Update(ctx, villaRule); err != nil {
			return err
		}

		updatedVillaRule = villaRule

		return nil
	}); err != nil {
		return dto.VillaRuleResponse{}, err
	}

	return mapper.ToVillaRuleResponse(updatedVillaRule), nil
}

func (s *VillaRuleService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		villaRule, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaRule == nil {
			return appError.NotFound(fmt.Sprintf("Villa rule not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, villaRule.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, villaRule)
	})
}
