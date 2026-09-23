package dto

import (
	"time"
)

type (
	VillaPackageResponse struct {
		ID          uint                       `json:"id"`
		Title       string                     `json:"title"`
		Subtitle    string                     `json:"subtitle"`
		MaxCapacity int                        `json:"max_capacity"`
		Price       float64                    `json:"price"`
		Lists       []VillaPackageListResponse `json:"lists"`
		SortOrder   int                        `json:"sort_order"`
		CreatedAt   time.Time                  `json:"createdAt"`
		UpdatedAt   time.Time                  `json:"updatedAt"`
	}

	VillaPackageRequest struct {
		Title       string  `json:"title" form:"title" binding:"required"`
		Subtitle    string  `json:"subtitle" form:"subtitle" binding:"required"`
		MaxCapacity int     `json:"max_capacity" form:"max_capacity" binding:"required,gt=0"`
		Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
		SortOrder   *int    `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)
