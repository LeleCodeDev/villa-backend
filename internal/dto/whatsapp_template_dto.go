package dto

import "time"

type (
	WhatsappTemplateQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	WhatsappTemplateRequest struct {
		Text string `json:"text" form:"text" binding:"required"`
	}

	WhatsappTemplateResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
