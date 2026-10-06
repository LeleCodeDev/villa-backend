package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"

	"github.com/lelecodedev/villa-backend/internal/image"
	"gorm.io/gorm"
)

type RentalOptionService struct {
	txManager *repository.TxManager
	repo      *repository.RentalOptionRepository
}

func NewRentalOptionService(
	txManager *repository.TxManager,
	repo *repository.RentalOptionRepository,
) *RentalOptionService {
	return &RentalOptionService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *RentalOptionService) GetAll(ctx context.Context, query dto.RentalOptionQuery) ([]dto.RentalOptionResponse, int64, error) {
	rentalOptions, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.RentalOptionResponse, 0, len(rentalOptions))
	for _, rentalOption := range rentalOptions {
		responses = append(responses, mapper.ToRentalOptionResponse(&rentalOption))
	}

	return responses, total, nil
}

func (s *RentalOptionService) GetByID(ctx context.Context, id uint) (dto.RentalOptionResponse, error) {
	rentalOption, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.RentalOptionResponse{}, err
	}
	if rentalOption == nil {
		return dto.RentalOptionResponse{}, appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", id))
	}

	return mapper.ToRentalOptionResponse(rentalOption), nil
}

func (s *RentalOptionService) Create(ctx context.Context, req dto.RentalOptionRequest) (dto.RentalOptionResponse, error) {
	var createdRentalOption *model.RentalOption

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

		var imagePath *string
		if req.Image != nil {
			path, err := image.SaveImage("uploads/rental_options", req.Image, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath = &path
		}

		rentalOption := mapper.ToRentalOptionModel(req, sortOrder, imagePath)
		if err := txRepo.Create(ctx, rentalOption); err != nil {
			image.DeleteImage(*imagePath)
			return err
		}

		createdRentalOption = rentalOption

		return nil
	}); err != nil {
		return dto.RentalOptionResponse{}, err
	}

	return mapper.ToRentalOptionResponse(createdRentalOption), nil
}

func (s *RentalOptionService) Update(ctx context.Context, req dto.RentalOptionRequest, id uint) (dto.RentalOptionResponse, error) {
	var updatedRentalOption *model.RentalOption

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		rentalOption, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if rentalOption == nil {
			return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", id))
		}

		oldOrder := rentalOption.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if newOrder != &oldOrder {
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

		imagepath := rentalOption.Image
		if req.Image != nil {
			if rentalOption.Image != nil {
				image.DeleteImage(*rentalOption.Image)
			}

			path, err := image.SaveImage("uploads/rental_options", req.Image, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath = &path
		}

		mapper.UpdateRentalOptionModel(rentalOption, req, *newOrder, imagepath)
		if err := txRepo.Update(ctx, rentalOption); err != nil {
			return err
		}

		updatedRentalOption = rentalOption

		return nil
	}); err != nil {
		return dto.RentalOptionResponse{}, err
	}

	return mapper.ToRentalOptionResponse(updatedRentalOption), nil
}

func (s *RentalOptionService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		rentalOption, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if rentalOption == nil {
			return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, rentalOption.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		if rentalOption.Image != nil {
			image.DeleteImage(*rentalOption.Image)
		}

		return txRepo.Delete(ctx, rentalOption)
	})
}
