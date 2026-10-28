package mapper

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToSeasonPeriodResponse(seasonPeriod *model.SeasonPeriod) dto.SeasonPeriodResponse {
	return dto.SeasonPeriodResponse{
		ID:        seasonPeriod.ID,
		Name:      seasonPeriod.Name,
		StartDate: seasonPeriod.StartDate,
		EndDate:   seasonPeriod.EndDate,
		CreatedAt: seasonPeriod.CreatedAt,
		UpdatedAt: seasonPeriod.UpdatedAt,
	}
}

func ToSeasonPeriodModel(req dto.SeasonPeriodRequest, startDate time.Time, endDate time.Time) *model.SeasonPeriod {
	return &model.SeasonPeriod{
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func UpdateSeasonPeriodModel(seasonPeriod *model.SeasonPeriod, req dto.SeasonPeriodRequest, startDate time.Time, endDate time.Time) {
	seasonPeriod.Name = req.Name
	seasonPeriod.StartDate = startDate
	seasonPeriod.EndDate = endDate
}
