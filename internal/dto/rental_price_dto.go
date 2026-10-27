package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	RentalPriceQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	RentalPriceResponse struct {
		ID             uint             `json:"id"`
		RentalOptionID uint             `json:"rental_option_id"`
		Season         model.SeasonType `json:"season"`
		DayType        model.DayType    `json:"day_type"`
		Price          float64          `json:"price"`
		CreatedAt      time.Time        `json:"created_at"`
		UpdatedAt      time.Time        `json:"updated_at"`
	}

	RentalPriceRequest struct {
		RentalOptionID uint             `json:"rental_option_id" form:"rental_option_id" binding:"required,gt=0"`
		Season         model.SeasonType `json:"season" form:"season" binding:"required,oneof=reguler high"`
		DayType        model.DayType    `json:"day_type" form:"day_type" binding:"required,oneof=weekday weekend"`
		Price          float64          `json:"price" form:"price" binding:"required,gt=0"`
	}
)

func (rpq *RentalPriceQuery) SetDefault() {
	rpq.setDefaultPagination()
	rpq.SetDefaultSort(SortAsc)
}
