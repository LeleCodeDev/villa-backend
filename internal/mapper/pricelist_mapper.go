package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToPricelistResponse(pricelist *model.Pricelist) dto.PricelistResponse {
	return dto.PricelistResponse{
		ID:               pricelist.ID,
		RentalOption:     ToRentalOptionResponse(&pricelist.RentalOption, nil),
		WeekdayPrice:     pricelist.WeekdayPrice,
		WeekendPrice:     pricelist.WeekdayPrice,
		LongWeekendPrice: pricelist.LongWeekendPrice,
		HighSeasonPrice:  pricelist.HighSeasonPrice,
		CreatedAt:        pricelist.CreatedAt,
		UpdatedAt:        pricelist.UpdatedAt,
	}
}

func ToPricelistModel(req dto.PricelistRequest, rentalOption *model.RentalOption) *model.Pricelist {
	return &model.Pricelist{
		RentalOption:     *rentalOption,
		RentalOptionID:   rentalOption.ID,
		WeekdayPrice:     req.WeekdayPrice,
		WeekendPrice:     req.WeekdayPrice,
		LongWeekendPrice: req.LongWeekendPrice,
		HighSeasonPrice:  req.HighSeasonPrice,
	}
}

func UpdatePricelistModel(pricelist *model.Pricelist, req dto.PricelistRequest, rentalOption *model.RentalOption) {
	pricelist.RentalOption = *rentalOption
	pricelist.RentalOptionID = rentalOption.ID
	pricelist.WeekdayPrice = req.WeekdayPrice
	pricelist.WeekendPrice = req.WeekdayPrice
	pricelist.LongWeekendPrice = req.LongWeekendPrice
	pricelist.HighSeasonPrice = req.HighSeasonPrice
}
