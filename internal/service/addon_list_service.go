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

type AddonListService struct {
	txManager *repository.TxManager
	repo      *repository.AddonListRepository
	addonRepo *repository.AddonRepository
}

func NewAddonListService(
	txManager *repository.TxManager,
	repo *repository.AddonListRepository,
	addonRepo *repository.AddonRepository,
) *AddonListService {
	return &AddonListService{
		txManager: txManager,
		repo:      repo,
		addonRepo: addonRepo,
	}
}

func (s *AddonListService) GetAll(ctx context.Context, query dto.AddonListQuery) ([]dto.AddonListResponse, int64, error) {
	addonLists, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.AddonListResponse, 0, len(addonLists))
	for _, addonList := range addonLists {
		responses = append(responses, mapper.ToAddonListResponse(&addonList))
	}

	return responses, total, nil
}

func (s *AddonListService) GetAllByAddonID(ctx context.Context, addonID uint, query dto.AddonListQuery) ([]dto.AddonListResponse, int64, error) {
	addonLists, total, err := s.repo.GetAllByAddonIDWithQuery(ctx, addonID, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.AddonListResponse, 0, len(addonLists))
	for _, addonList := range addonLists {
		responses = append(responses, mapper.ToAddonListResponse(&addonList))
	}

	return responses, total, nil
}

func (s *AddonListService) GetByID(ctx context.Context, id uint) (dto.AddonListResponse, error) {
	addonList, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.AddonListResponse{}, err
	}
	if addonList == nil {
		return dto.AddonListResponse{}, errors.NotFound(fmt.Sprintf("Addon list not found with ID: %d", id))
	}

	return mapper.ToAddonListResponse(addonList), nil
}

func (s *AddonListService) Create(ctx context.Context, req dto.AddonListRequest) (dto.AddonListResponse, error) {
	var createdAddonList *model.AddonList

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txAddonRepo := s.addonRepo.WithTx(tx)

		addon, err := txAddonRepo.GetByID(ctx, req.AddonID)
		if err != nil {
			return err
		}
		if addon == nil {
			return errors.NotFound(fmt.Sprintf("Addon not found with ID : %d", req.AddonID))
		}

		var sortOrder int
		if req.SortOrder == nil {
			maxOrder, err := txRepo.GetMaxSortOrder(ctx, addon.ID)
			if err != nil {
				return err
			}

			sortOrder = maxOrder + 1
		} else {
			sortOrder = *req.SortOrder
			if err := txRepo.ShiftSortOrder(ctx, addon.ID, sortOrder); err != nil {
				return err
			}
		}

		addonList := mapper.ToAddonListModel(req, *addon, sortOrder)
		if err := txRepo.Create(ctx, addonList); err != nil {
			return err
		}

		createdAddonList = addonList

		return nil
	}); err != nil {
		return dto.AddonListResponse{}, err
	}

	return mapper.ToAddonListResponse(createdAddonList), nil
}

func (s *AddonListService) Update(ctx context.Context, id uint, req dto.AddonListRequest) (dto.AddonListResponse, error) {
	var updatedAddonList *model.AddonList

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txAddonRepo := s.addonRepo.WithTx(tx)

		addonList, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if addonList == nil {
			return errors.NotFound(fmt.Sprintf("Addon list not found with ID: %d", id))
		}

		addon, err := txAddonRepo.GetByID(ctx, req.AddonID)
		if err != nil {
			return err
		}
		if addon == nil {
			return errors.NotFound(fmt.Sprintf("Addon not found with ID : %d", req.AddonID))
		}

		oldOrder := addonList.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if *newOrder != oldOrder {
			if *newOrder < oldOrder {
				if err := txRepo.UpdateSortOrderRange(ctx, addon.ID, *newOrder, oldOrder-1, +1); err != nil {
					return err
				}
			} else {
				if err := txRepo.UpdateSortOrderRange(ctx, addon.ID, oldOrder+1, *newOrder, -1); err != nil {
					return err
				}
			}
		}

		mapper.UpdateAddonListModel(addonList, *addon, req, *newOrder)
		if err := txRepo.Update(ctx, addonList); err != nil {
			return err
		}

		updatedAddonList = addonList

		return nil
	}); err != nil {
		return dto.AddonListResponse{}, err
	}

	return mapper.ToAddonListResponse(updatedAddonList), nil
}

func (s *AddonListService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		addonList, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if addonList == nil {
			return errors.NotFound(fmt.Sprintf("Addon list not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx, addonList.ID)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, addonList.AddonID, addonList.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, addonList)
	})
}
