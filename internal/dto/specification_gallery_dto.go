package dto

import (
	"mime/multipart"
	"time"
)

type (
	SpecificationGalleryResponse struct {
		ID        uint      `json:"id"`
		Image1    *string   `json:"image1"`
		Image2    *string   `json:"image2"`
		Image3    *string   `json:"image3"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	SpecificationGalleryRequest struct {
		Image1 *multipart.FileHeader `json:"image1" form:"image1" binding:"omitempty"`
		Image2 *multipart.FileHeader `json:"image2" form:"image2" binding:"omitempty"`
		Image3 *multipart.FileHeader `json:"image3" form:"image3" binding:"omitempty"`
	}
)
