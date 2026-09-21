package dto

import "time"

type (
	FaqResponse struct {
		ID        uint      `json:"id"`
		Question  string    `json:"question"`
		Answer    string    `json:"answer"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	FaqRequest struct {
		Question  string `json:"question" form:"question" binding:"required"`
		Answer    string `json:"answer" form:"answer" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)
