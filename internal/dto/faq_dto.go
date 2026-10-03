package dto

import (
	"time"
)

type (
	FaqQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	FaqResponse struct {
		ID        uint      `json:"id"`
		Question  string    `json:"question"`
		Answer    string    `json:"answer"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	FaqRequest struct {
		Question  string `json:"question" form:"question" binding:"required"`
		Answer    string `json:"answer" form:"answer" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (fq *FaqQuery) SetDefault() {
	fq.setDefaultPagination()
	fq.SetDefaultSort(SortAsc)
}
