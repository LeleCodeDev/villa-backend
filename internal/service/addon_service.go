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

type AddonService struct {
	txManager *repository.TxManager
	repo      *repository.AddonRepository
	listRepo  *repository.AddonListRepository
}

func NewAddonService(
	txManager *repository.TxManager,
	repo *repository.AddonRepository,
	listRepo *repository.AddonListRepository,
) *AddonService {
	return &AddonService{
		txManager: txManager,
		repo:      repo,
		listRepo:  listRepo,
	}
}

func (s *AddonService) GetAll(ctx context.Context, query dto.AddonQuery) ([]dto.AddonResponse, int64, error) {
	addons, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	addonIds := make([]uint, 0, len(addons))
	for _, addon := range addons {
		addonIds = append(addonIds, addon.ID)
	}

	addonLists, err := s.listRepo.GetAllByAddonIDIn(ctx, addonIds)
	if err != nil {
		return nil, 0, err
	}

	listMap := make(map[uint][]model.AddonList, len(addons))
	for _, list := range addonLists {
		listMap[list.AddonID] = append(listMap[list.AddonID], list)
	}

	responses := make([]dto.AddonResponse, 0, len(addons))
	for _, addon := range addons {
		responses = append(responses, mapper.ToAddonResponse(&addon, listMap[addon.ID]))
	}

	return responses, total, nil
}

func (s *AddonService) GetByID(ctx context.Context, id uint) (dto.AddonResponse, error) {
	addon, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.AddonResponse{}, err
	}
	if addon == nil {
		return dto.AddonResponse{}, appError.NotFound(fmt.Sprintf("Addon not found with ID : %d", id))
	}

	lists, err := s.listRepo.GetAllByAddonID(ctx, addon.ID)
	if err != nil {
		return dto.AddonResponse{}, err
	}

	return mapper.ToAddonResponse(addon, lists), nil
}

func (s *AddonService) Create(ctx context.Context, req dto.AddonRequest) (dto.AddonResponse, error) {
	var createdAddon *model.Addon

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

		addon := mapper.ToAddonModel(req, sortOrder)
		if err := txRepo.Create(ctx, addon); err != nil {
			return err
		}

		createdAddon = addon

		return nil
	}); err != nil {
		return dto.AddonResponse{}, err
	}

	return mapper.ToAddonResponse(createdAddon, nil), nil
}

func (s *AddonService) Update(ctx context.Context, id uint, req dto.AddonRequest) (dto.AddonResponse, error) {
	var updatedAddon *model.Addon

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		addon, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if addon == nil {
			return appError.NotFound(fmt.Sprintf("Addon not found with ID: %d", id))
		}

		oldOrder := addon.SortOrder
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

		mapper.UpdateAddonModel(addon, req, *newOrder)
		if err := txRepo.Update(ctx, addon); err != nil {
			return err
		}

		updatedAddon = addon

		return nil
	}); err != nil {
		return dto.AddonResponse{}, err
	}

	return mapper.ToAddonResponse(updatedAddon, nil), nil
}

func (s *AddonService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txListRepo := s.listRepo.WithTx(tx)

		addon, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if addon == nil {
			return appError.NotFound(fmt.Sprintf("Addon not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, addon.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		if err := txListRepo.DeleteByAddonID(ctx, id); err != nil {
			return err
		}

		return txRepo.Delete(ctx, addon)
	})
}
