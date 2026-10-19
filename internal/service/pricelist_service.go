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

type PricelistService struct {
	txManager        *repository.TxManager
	repo             *repository.PricelistRepository
	rentalOptionRepo *repository.RentalOptionRepository
}

func NewPricelistService(
	txManager *repository.TxManager,
	repo *repository.PricelistRepository,
	rentalOptionRepo *repository.RentalOptionRepository,
) *PricelistService {
	return &PricelistService{
		txManager:        txManager,
		repo:             repo,
		rentalOptionRepo: rentalOptionRepo,
	}
}

func (s *PricelistService) GetAll(ctx context.Context, query dto.PricelistQuery) ([]dto.PricelistResponse, int64, error) {
	pricelists, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PricelistResponse, 0, len(pricelists))
	for _, pricelist := range pricelists {
		responses = append(responses, mapper.ToPricelistResponse(&pricelist))
	}

	return responses, total, nil
}

func (s *PricelistService) GetByID(ctx context.Context, id uint) (dto.PricelistResponse, error) {
	pricelist, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.PricelistResponse{}, err
	}
	if pricelist == nil {
		return dto.PricelistResponse{}, appError.NotFound(fmt.Sprintf("Pricelist not found with ID: %d", id))
	}

	return mapper.ToPricelistResponse(pricelist), nil
}

func (s *PricelistService) Create(ctx context.Context, req dto.PricelistRequest) (dto.PricelistResponse, error) {
	var createdPricelist *model.Pricelist

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txRentalOptionRepo := s.rentalOptionRepo.WithTx(tx)

		rentalOption, err := txRentalOptionRepo.GetByID(ctx, req.RentalOptionID)
		if err != nil {
			return err
		}
		if rentalOption == nil {
			return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", req.RentalOptionID))
		}

		pricelist := mapper.ToPricelistModel(req, rentalOption)
		if err := txRepo.Create(ctx, pricelist); err != nil {
			return err
		}

		createdPricelist = pricelist

		return nil
	}); err != nil {
		return dto.PricelistResponse{}, err
	}

	return mapper.ToPricelistResponse(createdPricelist), nil
}

func (s *PricelistService) Update(ctx context.Context, id uint, req dto.PricelistRequest) (dto.PricelistResponse, error) {
	var updatedPricelist *model.Pricelist

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txRentalOptionRepo := s.rentalOptionRepo.WithTx(tx)

		pricelist, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if pricelist == nil {
			return appError.NotFound(fmt.Sprintf("Pricelist not found with ID: %d", id))
		}

		rentalOption := &pricelist.RentalOption
		if pricelist.RentalOptionID != req.RentalOptionID {

			rentalOption, err = txRentalOptionRepo.GetByID(ctx, req.RentalOptionID)
			if err != nil {
				return err
			}
			if rentalOption == nil {
				return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", req.RentalOptionID))
			}
		}

		mapper.UpdatePricelistModel(pricelist, req, rentalOption)
		if err := txRepo.Update(ctx, pricelist); err != nil {
			return err
		}

		updatedPricelist = pricelist

		return nil
	}); err != nil {
		return dto.PricelistResponse{}, err
	}

	return mapper.ToPricelistResponse(updatedPricelist), nil
}

func (s *PricelistService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		pricelist, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if pricelist == nil {
			return appError.NotFound(fmt.Sprintf("Pricelist not found with ID: %d", id))
		}

		return txRepo.Delete(ctx, pricelist)
	})
}
