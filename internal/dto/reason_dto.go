package dto

import "time"

type (
	ReasonRequest struct {
		Text string `json:"text" form:"text" binding:"required"`
	}

	ReasonResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
