package dto

import (
	"mime/multipart"
	"time"
)

type (
	ApplicationResponse struct {
		ID            uint      `json:"id"`
		Title         string    `json:"title"`
		Logo          *string   `json:"logo"`
		PhoneNumber   string    `json:"phone_number"`
		Email         *string   `json:"email"`
		InstagramLink *string   `json:"instagram_link"`
		TikTokLink    *string   `json:"tiktok_link"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	ApplicationRequest struct {
		Title         string                `json:"title" form:"title" binding:"required"`
		Logo          *multipart.FileHeader `form:"logo" binding:"omitempty"`
		PhoneNumber   string                `json:"phone_number" form:"phone_number" binding:"required"`
		Email         *string               `json:"email" form:"email" binding:"omitempty,email"`
		InstagramLink *string               `json:"instagram_link" form:"instagram_link" binding:"omitempty,http_url"`
		TikTokLink    *string               `json:"tiktok_link" form:"tiktok_link" binding:"omitempty,http_url"`
	}
)
