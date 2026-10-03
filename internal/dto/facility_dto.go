package dto

import (
	"time"
)

type (
	FacilityQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	FacilityResponse struct {
		ID          uint      `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Logo        string    `json:"logo"`
		SortOrder   int       `json:"sort_order"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	FacilityRequest struct {
		Title       string `json:"title" form:"title" binding:"required"`
		Description string `json:"description" form:"description" binding:"required"`
		Logo        string `json:"logo" form:"logo" binding:"required"`
		SortOrder   *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (fq *FacilityQuery) SetDefault() {
	fq.setDefaultPagination()
	fq.SetDefaultSort(SortAsc)
}
