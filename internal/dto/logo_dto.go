package dto

import (
	"mime/multipart"
	"time"
)

type (
	LogoQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	LogoResponse struct {
		ID        uint      `json:"id"`
		Image     string    `json:"image"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	LogoRequest struct {
		Image *multipart.FileHeader `json:"image" form:"image" binding:"required"`
	}
)

func (lq *LogoQuery) SetDefault() {
	lq.setDefaultPagination()
	lq.SetDefaultSort(SortAsc)
}
