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
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	SpecificationGalleryRequest struct {
		Image1 *multipart.FileHeader `form:"image1" binding:"omitempty"`
		Image2 *multipart.FileHeader `form:"image2" binding:"omitempty"`
		Image3 *multipart.FileHeader `form:"image3" binding:"omitempty"`
	}
)
