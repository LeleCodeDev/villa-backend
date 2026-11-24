package dto

import (
	"mime/multipart"
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	GalleryQuery struct {
		PaginationQuery
		SortQuery
		Type   model.GalleryType `form:"type" binding:"omitempty,oneof=booking wedding"`
		Unpage bool              `form:"unpage"`
	}

	GalleryResponse struct {
		ID          uint              `json:"id"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Image       *string           `json:"image"`
		Type        model.GalleryType `json:"type"`
		SortOrder   int               `json:"sort_order"`
		CreatedAt   time.Time         `json:"created_at"`
		UpdatedAt   time.Time         `json:"updated_at"`
	}

	GalleryRequest struct {
		Title       string                `json:"title" form:"title" binding:"required"`
		Description string                `json:"description" form:"description" binding:"required"`
		SortOrder   *int                  `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
		Image       *multipart.FileHeader `json:"image" form:"image" binding:"omitempty"`
		Type        model.GalleryType     `json:"type" form:"type" binding:"required,oneof=booking wedding"`
	}
)

func (gq *GalleryQuery) SetDefault() {
	gq.setDefaultPagination()
	gq.SetDefaultSort(SortAsc)
}
