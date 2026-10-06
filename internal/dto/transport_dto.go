package dto

import "time"

type (
	TransportQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	TransportResponse struct {
		ID           uint      `json:"id"`
		Vehicle      string    `json:"vehicle"`
		Distance     string    `json:"distance"`
		TimeRequired string    `json:"time_required"`
		SortOrder    int       `json:"sort_order"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	TransportRequest struct {
		AttractionID uint   `json:"attraction_id" binding:"required,gt=0"`
		Vehicle      string `json:"vehicle" form:"vehicle" binding:"required"`
		Distance     string `json:"distance" form:"distance" binding:"required"`
		TimeRequired string `json:"time_required" form:"time_required" binding:"required"`
		SortOrder    *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (tq *TransportQuery) SetDefault() {
	tq.setDefaultPagination()
	tq.SetDefaultSort(SortAsc)
}
