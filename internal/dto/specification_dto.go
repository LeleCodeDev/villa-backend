package dto

import (
	"time"
)

type (
	SpecificationQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	SpecificationResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		Logo      string    `json:"logo"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	SpecificationRequest struct {
		Text      string `json:"text" form:"text" binding:"required"`
		Logo      string `json:"logo" form:"logo" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (sq *SpecificationQuery) SetDefault() {
	sq.setDefaultPagination()
	sq.SetDefaultSort(SortAsc)
}
