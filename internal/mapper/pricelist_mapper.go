package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToPricelistResponse(pricelist *model.Pricelist) dto.PricelistResponse {
	return dto.PricelistResponse{
		ID:               pricelist.ID,
		Title:            pricelist.Title,
		Description:      pricelist.Description,
		MinCapacity:      pricelist.MinCapacity,
		MaxCapacity:      pricelist.MaxCapacity,
		WeekdayPrice:     pricelist.WeekdayPrice,
		WeekendPrice:     pricelist.WeekdayPrice,
		LongWeekendPrice: pricelist.LongWeekendPrice,
		HighSeasonPrice:  pricelist.HighSeasonPrice,
		CreatedAt:        pricelist.CreatedAt,
		UpdatedAt:        pricelist.UpdatedAt,
	}
}

func ToPricelistModel(req dto.PricelistRequest) *model.Pricelist {
	return &model.Pricelist{
		Title:            req.Title,
		Description:      req.Description,
		MinCapacity:      req.MinCapacity,
		MaxCapacity:      req.MaxCapacity,
		WeekdayPrice:     req.WeekdayPrice,
		WeekendPrice:     req.WeekdayPrice,
		LongWeekendPrice: req.LongWeekendPrice,
		HighSeasonPrice:  req.HighSeasonPrice,
	}
}

func UpdatePricelistModel(pricelist *model.Pricelist, req dto.PricelistRequest) {
	pricelist.Title = req.Title
	pricelist.Description = req.Description
	pricelist.MinCapacity = req.MinCapacity
	pricelist.MaxCapacity = req.MaxCapacity
	pricelist.WeekdayPrice = req.WeekdayPrice
	pricelist.WeekendPrice = req.WeekdayPrice
	pricelist.LongWeekendPrice = req.LongWeekendPrice
	pricelist.HighSeasonPrice = req.HighSeasonPrice
}
