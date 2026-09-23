package dto

import "time"

type (
	TestimonyResponse struct {
		ID        uint      `json:"id"`
		Star      int       `json:"star"`
		Comment   string    `json:"comment"`
		Username  string    `json:"username"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	TestimonyRequest struct {
		Star      int    `json:"star" form:"star" binding:"required,gte=0,lte=5"`
		Comment   string `json:"comment" form:"comment" binding:"required"`
		Username  string `json:"username" form:"username" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)
