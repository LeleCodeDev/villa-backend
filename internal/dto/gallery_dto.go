package dto

import (
	"mime/multipart"
	"time"
)

type (
	GalleryResponse struct {
		ID          uint      `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Image       *string   `json:"image"`
		SortOrder   int       `json:"sort_order"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	GalleryRequest struct {
		Title       string                `json:"title" form:"title" binding:"required"`
		Description string                `json:"description" form:"description" binding:"required"`
		SortOrder   *int                  `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
		Image       *multipart.FileHeader `form:"image" binding:"omitempty"`
	}
)
