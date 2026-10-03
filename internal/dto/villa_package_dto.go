package dto

import (
	"time"
)

type (
	VillaPackageQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	VillaPackageResponse struct {
		ID          uint                       `json:"id"`
		Title       string                     `json:"title"`
		Subtitle    string                     `json:"subtitle"`
		MaxCapacity int                        `json:"max_capacity"`
		Price       float64                    `json:"price"`
		Lists       []VillaPackageListResponse `json:"lists"`
		SortOrder   int                        `json:"sort_order"`
		CreatedAt   time.Time                  `json:"created_at"`
		UpdatedAt   time.Time                  `json:"updated_at"`
	}

	VillaPackageRequest struct {
		Title       string  `json:"title" form:"title" binding:"required"`
		Subtitle    string  `json:"subtitle" form:"subtitle" binding:"required"`
		MaxCapacity int     `json:"max_capacity" form:"max_capacity" binding:"required,gt=0"`
		Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
		SortOrder   *int    `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (vpq *VillaPackageQuery) SetDefault() {
	vpq.setDefaultPagination()
	vpq.SetDefaultSort(SortAsc)
}
