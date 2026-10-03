package dto

import "time"

type (
	PricelistRequest struct {
		Title            string  `json:"title" form:"title" binding:"required"`
		Description      string  `json:"description" form:"description" binding:"required"`
		MinCapacity      int     `json:"min_capacity" form:"min_capacity" binding:"required,gt=0"`
		MaxCapacity      int     `json:"max_capacity" form:"max_capacity" binding:"required,gt=0"`
		WeekdayPrice     float64 `json:"weekday_price" form:"weekday_price" binding:"required,gt=0"`
		WeekendPrice     float64 `json:"weekend_price" form:"weekend_price" binding:"required,gt=0"`
		LongWeekendPrice float64 `json:"long_weekend_price" form:"long_weekend_price" binding:"required,gt=0"`
		HighSeasonPrice  float64 `json:"high_season_price" form:"high_season_price" binding:"required,gt=0"`
	}

	PricelistResponse struct {
		ID               uint      `json:"id"`
		Title            string    `json:"title"`
		Description      string    `json:"description"`
		MinCapacity      int       `json:"min_capacity"`
		MaxCapacity      int       `json:"max_capacity"`
		WeekdayPrice     float64   `json:"weekday_price"`
		WeekendPrice     float64   `json:"weekend_price"`
		LongWeekendPrice float64   `json:"long_weekend_price"`
		HighSeasonPrice  float64   `json:"high_season_price"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
	}
)
