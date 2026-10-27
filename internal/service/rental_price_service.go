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

type RentalPriceService struct {
	txManager        *repository.TxManager
	repo             *repository.RentalPriceRepository
	rentalOptionRepo *repository.RentalOptionRepository
}

func NewRentalPriceService(
	txManager *repository.TxManager,
	repo *repository.RentalPriceRepository,
	rentalOptionRepo *repository.RentalOptionRepository,
) *RentalPriceService {
	return &RentalPriceService{
		txManager:        txManager,
		repo:             repo,
		rentalOptionRepo: rentalOptionRepo,
	}
}

func (s *RentalPriceService) GetAll(ctx context.Context, query dto.RentalPriceQuery) ([]dto.RentalPriceResponse, int64, error) {
	rentalPrices, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.RentalPriceResponse, 0, len(rentalPrices))
	for _, rentalPrice := range rentalPrices {
		responses = append(responses, mapper.ToRentalPriceResponse(&rentalPrice))
	}

	return responses, total, nil
}

func (s *RentalPriceService) GetAllByRentalOptionID(ctx context.Context, rentalOptionID uint) ([]dto.RentalPriceResponse, error) {
	rentalPrices, err := s.repo.GetAllByRentalOptionID(ctx, rentalOptionID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RentalPriceResponse, 0, len(rentalPrices))
	for _, rentalPrice := range rentalPrices {
		responses = append(responses, mapper.ToRentalPriceResponse(&rentalPrice))
	}

	return responses, nil
}

func (s *RentalPriceService) GetByID(ctx context.Context, id uint) (dto.RentalPriceResponse, error) {
	rentalPrice, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.RentalPriceResponse{}, err
	}
	if rentalPrice == nil {
		return dto.RentalPriceResponse{}, appError.NotFound(fmt.Sprintf("Rental price not found with ID: %d", id))
	}

	return mapper.ToRentalPriceResponse(rentalPrice), nil
}

func (s *RentalPriceService) Create(ctx context.Context, req dto.RentalPriceRequest) (dto.RentalPriceResponse, error) {
	var createdRentalPrice *model.RentalPrice

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

		exist, err := txRepo.ExistByRentalOptionIDAndSeasonAndDayType(ctx, req.RentalOptionID, req.Season, req.DayType)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist(fmt.Sprintf("Rental price for %s %s already exist for rental option with ID: %d", req.Season, req.DayType, req.RentalOptionID))
		}

		rentalPrice := mapper.ToRentalPriceModel(req, rentalOption)
		if err := txRepo.Create(ctx, rentalPrice); err != nil {
			return err
		}

		createdRentalPrice = rentalPrice

		return nil
	}); err != nil {
		return dto.RentalPriceResponse{}, err
	}

	return mapper.ToRentalPriceResponse(createdRentalPrice), nil
}

func (s *RentalPriceService) Update(ctx context.Context, id uint, req dto.RentalPriceRequest) (dto.RentalPriceResponse, error) {
	var updatedRentalPrice *model.RentalPrice

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txRentalOptionRepo := s.rentalOptionRepo.WithTx(tx)

		rentalPrice, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if rentalPrice == nil {
			return appError.NotFound(fmt.Sprintf("Rental price not found with ID: %d", id))
		}

		rentalOption := &rentalPrice.RentalOption
		isRentalOptionChanged := req.RentalOptionID != rentalPrice.RentalOptionID

		if isRentalOptionChanged {
			rentalOption, err = txRentalOptionRepo.GetByID(ctx, req.RentalOptionID)
			if err != nil {
				return err
			}
			if rentalOption == nil {
				return appError.NotFound(fmt.Sprintf("Rental option not found with ID: %d", req.RentalOptionID))
			}
		}

		isSeasonChanged := req.Season != rentalPrice.Season
		isDayTypeChanged := req.DayType != rentalPrice.DayType

		if isRentalOptionChanged || isSeasonChanged || isDayTypeChanged {
			exist, err := txRepo.ExistByRentalOptionIDAndSeasonAndDayType(ctx, req.RentalOptionID, req.Season, req.DayType)
			if err != nil {
				return err
			}
			if exist {
				return appError.AlreadyExist(fmt.Sprintf("Rental price for %s %s already exist for rental option with ID: %d", req.Season, req.DayType, req.RentalOptionID))
			}
		}

		mapper.UpdateRentalPriceModel(rentalPrice, req, rentalOption)
		if err := txRepo.Update(ctx, rentalPrice); err != nil {
			return err
		}

		updatedRentalPrice = rentalPrice

		return nil
	}); err != nil {
		return dto.RentalPriceResponse{}, err
	}

	return mapper.ToRentalPriceResponse(updatedRentalPrice), nil
}
