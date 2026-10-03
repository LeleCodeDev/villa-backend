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
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	ApplicationRequest struct {
		Title       string                `json:"title" form:"title" binding:"required"`
		Logo        *multipart.FileHeader `form:"logo" binding:"omitempty"`
		PhoneNumber string                `json:"phone_number" form:"phone_number" binding:"required"`
	}
)
