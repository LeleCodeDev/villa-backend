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

type VillaPackageService struct {
	txManager *repository.TxManager
	repo      *repository.VillaPackageRepository
	listRepo  *repository.VillaPackageListRepository
}

func NewVillaPackageService(
	txManager *repository.TxManager,
	repo *repository.VillaPackageRepository,
	listRepo *repository.VillaPackageListRepository,
) *VillaPackageService {
	return &VillaPackageService{
		txManager: txManager,
		repo:      repo,
		listRepo:  listRepo,
	}
}

func (s *VillaPackageService) GetAll(ctx context.Context) ([]dto.VillaPackageResponse, error) {
	villaPackages, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	villaPackageIds := make([]uint, 0, len(villaPackages))
	for _, villaPackage := range villaPackages {
		villaPackageIds = append(villaPackageIds, villaPackage.ID)
	}

	allLists, err := s.listRepo.GetAllByVillaPackageIDIn(ctx, villaPackageIds)
	if err != nil {
		return nil, err
	}

	listMap := make(map[uint][]dto.VillaPackageListResponse, len(villaPackages))
	for _, list := range allLists {
		listMap[list.VillaPackageID] = append(
			listMap[list.VillaPackageID],
			mapper.ToVillaPackageListResponse(&list),
		)
	}

	responses := make([]dto.VillaPackageResponse, 0, len(villaPackages))
	for _, villaPackage := range villaPackages {
		responses = append(
			responses,
			mapper.ToVillaPackageResponse(&villaPackage,
				listMap[villaPackage.ID]),
		)
	}

	return responses, nil
}

func (s *VillaPackageService) GetByID(ctx context.Context, id uint) (dto.VillaPackageResponse, error) {
	villaPackage, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.VillaPackageResponse{}, err
	}
	if villaPackage == nil {
		return dto.VillaPackageResponse{}, errors.NotFound(fmt.Sprintf("Villa package not found with ID: %d", id))
	}

	lists, err := s.listRepo.GetAllByVillaPackageID(ctx, villaPackage.ID)
	if err != nil {
		return dto.VillaPackageResponse{}, err
	}

	listsResponse := make([]dto.VillaPackageListResponse, 0, len(lists))
	for _, list := range lists {
		listsResponse = append(listsResponse, mapper.ToVillaPackageListResponse(&list))
	}

	return mapper.ToVillaPackageResponse(villaPackage, listsResponse), nil
}

func (s *VillaPackageService) Create(ctx context.Context, req dto.VillaPackageRequest) (dto.VillaPackageResponse, error) {
	var createdVillaPackage *model.VillaPackage

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

		villaPackage := mapper.ToVillaPackageModel(req, sortOrder)
		if err := txRepo.Create(ctx, villaPackage); err != nil {
			return err
		}

		createdVillaPackage = villaPackage

		return nil
	}); err != nil {
		return dto.VillaPackageResponse{}, err
	}

	return mapper.ToVillaPackageResponse(createdVillaPackage, nil), nil
}

func (s *VillaPackageService) Update(ctx context.Context, id uint, req dto.VillaPackageRequest) (dto.VillaPackageResponse, error) {
	var updatedVillaPackage *model.VillaPackage

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		villaPackage, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaPackage == nil {
			return errors.NotFound(fmt.Sprintf("Villa package not found with ID: %d", id))
		}

		oldOrder := villaPackage.SortOrder
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

		mapper.UpdateVillaPackageModel(villaPackage, req, *newOrder)
		if err := txRepo.Update(ctx, villaPackage); err != nil {
			return err
		}

		updatedVillaPackage = villaPackage

		return nil
	}); err != nil {
		return dto.VillaPackageResponse{}, err
	}

	return mapper.ToVillaPackageResponse(updatedVillaPackage, nil), nil
}

func (s *VillaPackageService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txListRepo := s.listRepo.WithTx(tx)

		villaPackage, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if villaPackage == nil {
			return errors.NotFound(fmt.Sprintf("Villa package not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, villaPackage.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		if err := txListRepo.DeleteByVillaPackageID(ctx, id); err != nil {
			return err
		}

		return txRepo.Delete(ctx, villaPackage)
	})
}
