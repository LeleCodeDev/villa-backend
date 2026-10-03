package dto

import (
	"time"
)

type (
	AddonQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	AddonResponse struct {
		ID          uint                `json:"id"`
		Title       string              `json:"title"`
		Subtitle    string              `json:"subtitle"`
		Description string              `json:"description"`
		Price       float64             `json:"price"`
		Lists       []AddonListResponse `json:"lists"`
		SortOrder   int                 `json:"sort_order"`
		CreatedAt   time.Time           `json:"created_at"`
		UpdatedAt   time.Time           `json:"updated_at"`
	}

	AddonRequest struct {
		Title       string  `json:"title" form:"title" binding:"required"`
		Subtitle    string  `json:"subtitle" form:"subtitle" binding:"required"`
		Description string  `json:"description" form:"description" binding:"required"`
		Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
		SortOrder   *int    `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (aq *AddonQuery) SetDefault() {
	aq.setDefaultPagination()
	aq.SetDefaultSort(SortAsc)
}
