package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToRentalPriceResponse(rentalPrice *model.RentalPrice) dto.RentalPriceResponse {
	return dto.RentalPriceResponse{
		ID:             rentalPrice.ID,
		RentalOptionID: rentalPrice.RentalOptionID,
		Season:         rentalPrice.Season,
		DayType:        rentalPrice.DayType,
		Price:          rentalPrice.Price,
		CreatedAt:      rentalPrice.CreatedAt,
		UpdatedAt:      rentalPrice.UpdatedAt,
	}
}

func ToRentalPriceModel(req dto.RentalPriceRequest, rentalOption *model.RentalOption) *model.RentalPrice {
	return &model.RentalPrice{
		RentalOptionID: rentalOption.ID,
		RentalOption:   *rentalOption,
		Season:         req.Season,
		DayType:        req.DayType,
		Price:          req.Price,
	}
}

func UpdateRentalPriceModel(rentalPrice *model.RentalPrice, req dto.RentalPriceRequest, rentalOption *model.RentalOption) {
	rentalPrice.RentalOptionID = rentalOption.ID
	rentalPrice.RentalOption = *rentalOption
	rentalPrice.Season = req.Season
	rentalPrice.DayType = req.DayType
	rentalPrice.Price = req.Price
}
