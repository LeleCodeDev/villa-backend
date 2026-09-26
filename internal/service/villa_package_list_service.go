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

type VillaPackageListService struct {
	txManager        *repository.TxManager
	repo             *repository.VillaPackageListRepository
	villaPackageRepo *repository.VillaPackageRepository
}

func NewVillaPackageListService(
	txManager *repository.TxManager,
	repo *repository.VillaPackageListRepository,
	villaPackageRepo *repository.VillaPackageRepository,
) *VillaPackageListService {
	return &VillaPackageListService{
		txManager:        txManager,
		repo:             repo,
		villaPackageRepo: villaPackageRepo,
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

func (s *VillaPackageListService) GetByID(ctx context.Context, id uint) (dto.VillaPackageListResponse, error) {
	villaPackageList, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.VillaPackageListResponse{}, err
	}
	if villaPackageList == nil {
		return dto.VillaPackageListResponse{}, errors.NotFound(fmt.Sprintf("Villa package list not found with ID: %d", id))
	}

	return mapper.ToVillaPackageListResponse(villaPackageList), nil
}

func (s *VillaPackageListService) Create(ctx context.Context, req dto.VillaPackageListRequest) (dto.VillaPackageListResponse, error) {
	var createdVillaPackageList *model.VillaPackageList

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txVillaPackageRepo := s.villaPackageRepo.WithTx(tx)

		villaPackage, err := txVillaPackageRepo.GetByID(ctx, req.VillaPackageID)
		if err != nil {
			return err
		}
		if villaPackage == nil {
			return errors.NotFound(fmt.Sprintf("Villa package not found with ID : %d", req.VillaPackageID))
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

		villaPackageList := mapper.ToVillaPackageListModel(req, *villaPackage, sortOrder)
		if err := txRepo.Create(ctx, villaPackageList); err != nil {
			return err
		}

		createdVillaPackageList = villaPackageList

		return nil
	}); err != nil {
		return dto.VillaPackageListResponse{}, err
	}

	return mapper.ToVillaPackageListResponse(createdVillaPackageList), nil
}

func (s *VillaPackageListService) Update(ctx context.Context, id uint, req dto.VillaPackageListRequest) (dto.VillaPackageListResponse, error) {
	var updatedVillaPackageList *model.VillaPackageList

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txVillaPackageRepo := s.villaPackageRepo.WithTx(tx)

		villaPackageList, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaPackageList == nil {
			return errors.NotFound(fmt.Sprintf("Villa package list not found with ID: %d", id))
		}

		villaPackage, err := txVillaPackageRepo.GetByID(ctx, req.VillaPackageID)
		if err != nil {
			return err
		}
		if villaPackage == nil {
			return errors.NotFound(fmt.Sprintf("Villa package not found with ID : %d", req.VillaPackageID))
		}

		oldOrder := villaPackageList.SortOrder
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

		mapper.UpdateVillaPackageListModel(villaPackageList, *villaPackage, req, *newOrder)
		if err := txRepo.Update(ctx, villaPackageList); err != nil {
			return err
		}

		updatedVillaPackageList = villaPackageList

		return nil
	}); err != nil {
		return dto.VillaPackageListResponse{}, err
	}

	return mapper.ToVillaPackageListResponse(updatedVillaPackageList), nil
}

func (s *VillaPackageListService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		villaPackageList, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaPackageList == nil {
			return errors.NotFound(fmt.Sprintf("Villa package list not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, villaPackageList.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, villaPackageList)
	})
}
