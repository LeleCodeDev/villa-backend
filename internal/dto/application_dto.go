package dto

import (
	"mime/multipart"
	"time"
)

type (
	ApplicationResponse struct {
		ID          uint      `json:"id"`
		Title       string    `json:"title"`
		Logo        *string   `json:"logo"`
		PhoneNumber string    `json:"phone_number"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
	}

	ApplicationRequest struct {
		Title       string                `json:"title" form:"title" binding:"required"`
		Logo        *multipart.FileHeader `form:"image" binding:"omitempty"`
		PhoneNumber string                `json:"phone_number" form:"phone_number" binding:"required"`
	}
)
