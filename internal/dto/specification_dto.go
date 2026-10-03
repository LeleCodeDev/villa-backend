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
		ID        uint         `json:"id"`
		Text      string       `json:"text"`
		Logo      LogoResponse `json:"logo"`
		SortOrder int          `json:"sort_order"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedAt time.Time    `json:"updated_at"`
	}

	SpecificationRequest struct {
		Text      string `json:"text" form:"text" binding:"required"`
		LogoID    uint   `json:"logo_id" binding:"required,gt=0"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (sq *SpecificationQuery) SetDefault() {
	sq.setDefaultPagination()
	sq.SetDefaultSort(SortAsc)
}
