// Package dto
package dto

import "time"

type (
	HeroSectionResponse struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title"`
		Subtitle  string    `json:"subtitle"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	HeroSectionRequest struct {
		Title    string `json:"title" form:"title" binding:"required"`
		Subtitle string `json:"subtitle" form:"subtitle" binding:"required"`
	}
)
