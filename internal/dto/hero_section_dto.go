package dto

import "time"

type (
	HeroSectionResponse struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Subtitle  string    `json:"subtitle"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	HeroSectionRequest struct {
		Title    string `json:"title" form:"title" binding:"required"`
		Subtitle string `json:"subtitle" form:"subtitle" binding:"required"`
	}
)
