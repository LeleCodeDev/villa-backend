package dto

import (
	"mime/multipart"
	"time"
)

type (
	AttractionQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	AttractionResponse struct {
		ID         uint                `json:"id"`
		Title      string              `json:"title"`
		Subtitle   string              `json:"subtitle"`
		Price      float64             `json:"price"`
		Transports []TransportResponse `json:"transports"`
		SortOrder  int                 `json:"sort_order"`
		Image      *string             `json:"image"`
		CreatedAt  time.Time           `json:"created_at"`
		UpdatedAt  time.Time           `json:"updated_at"`
	}

	AttractionRequest struct {
		Title     string                `json:"title" form:"title" binding:"required"`
		Subtitle  string                `json:"subtitle" form:"subtitle" binding:"required"`
		Price     float64               `json:"price" form:"price" binding:"required,gt=0"`
		SortOrder *int                  `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
		Image     *multipart.FileHeader `json:"image" form:"image" binding:"omitempty"`
	}
)

func (aq *AttractionQuery) SetDefault() {
	aq.setDefaultPagination()
	aq.SetDefaultSort(SortAsc)
}
