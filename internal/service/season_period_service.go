package service

import (
	"context"
	"fmt"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type SeasonPeriodService struct {
	txManager *repository.TxManager
	repo      *repository.SeasonPeriodRepository
}

func NewSeasionPeriodService(
	txManager *repository.TxManager,
	repo *repository.SeasonPeriodRepository,
) *SeasonPeriodService {
	return &SeasonPeriodService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *SeasonPeriodService) GetAll(ctx context.Context, query dto.SeasonPeriodQuery) ([]dto.SeasonPeriodResponse, int64, error) {
	seasonPeriods, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.SeasonPeriodResponse, 0, len(seasonPeriods))
	for _, seasonPeriod := range seasonPeriods {
		responses = append(responses, mapper.ToSeasonPeriodResponse(&seasonPeriod))
	}

	return responses, total, nil
}

func (s *SeasonPeriodService) GetByID(ctx context.Context, id uint) (dto.SeasonPeriodResponse, error) {
	seasonPeriod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.SeasonPeriodResponse{}, err
	}
	if seasonPeriod == nil {
		return dto.SeasonPeriodResponse{}, appError.NotFound(fmt.Sprintf("Season period not found with ID: %d", id))
	}

	return mapper.ToSeasonPeriodResponse(seasonPeriod), nil
}

func parseDates(startDate string, endDate string) (time.Time, time.Time, error) {
	loc := config.Env.Location

	start, err := time.ParseInLocation("2006-01-02", startDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid start_date, expected YYYY-MM-DD")
	}

	end, err := time.ParseInLocation("2006-01-02", endDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, appError.BadRequest("invalid end_date, expected YYYY-MM-DD")
	}

	if end.Before(start) {
		return time.Time{}, time.Time{}, appError.BadRequest("end_date must be on or after start_date")
	}

	return start, end, nil
}

func (s *SeasonPeriodService) Create(ctx context.Context, req dto.SeasonPeriodRequest) (dto.SeasonPeriodResponse, error) {
	start, end, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return dto.SeasonPeriodResponse{}, err
	}

	var createdSeasonPeriod *model.SeasonPeriod

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		conflicts, err := txRepo.GetSeasonConflict(ctx, start, end, 0)
		if err != nil {
			return err
		}

		if len(conflicts) > 0 {
			conflictResponses := make([]dto.SeasonPeriodResponse, 0, len(conflicts))
			for _, c := range conflicts {
				conflictResponses = append(conflictResponses, mapper.ToSeasonPeriodResponse(&c))
			}

			return appError.Conflict("Dates overlap with existing season periods").WithErrors(conflictResponses)
		}

		seasonPeriod := mapper.ToSeasonPeriodModel(req, start, end)
		if err := txRepo.Create(ctx, seasonPeriod); err != nil {
			return err
		}

		createdSeasonPeriod = seasonPeriod

		return nil
	}); err != nil {
		return dto.SeasonPeriodResponse{}, err
	}

	return mapper.ToSeasonPeriodResponse(createdSeasonPeriod), nil
}

func (s *SeasonPeriodService) Update(ctx context.Context, id uint, req dto.SeasonPeriodRequest) (dto.SeasonPeriodResponse, error) {
	start, end, err := parseDates(req.StartDate, req.EndDate)
	if err != nil {
		return dto.SeasonPeriodResponse{}, err
	}

	var updatedSeasonPeriod *model.SeasonPeriod

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		seasonPeriod, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if seasonPeriod == nil {
			return appError.NotFound(fmt.Sprintf("Season period not found with ID: %d", id))
		}

		conflicts, err := txRepo.GetSeasonConflict(ctx, start, end, id)
		if err != nil {
			return err
		}

		if len(conflicts) > 0 {
			conflictResponses := make([]dto.SeasonPeriodResponse, 0, len(conflicts))
			for _, c := range conflicts {
				conflictResponses = append(conflictResponses, mapper.ToSeasonPeriodResponse(&c))
			}

			return appError.Conflict("Dates overlap with existing season periods").WithErrors(conflictResponses)
		}

		mapper.UpdateSeasonPeriodModel(seasonPeriod, req, start, end)
		if err := txRepo.Update(ctx, seasonPeriod); err != nil {
			return err
		}

		updatedSeasonPeriod = seasonPeriod

		return nil
	}); err != nil {
		return dto.SeasonPeriodResponse{}, err
	}

	return mapper.ToSeasonPeriodResponse(updatedSeasonPeriod), nil
}
