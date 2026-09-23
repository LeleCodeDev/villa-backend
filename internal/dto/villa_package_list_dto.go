package dto

import "time"

type (
	VillaPackageListResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	VillaPackageListRequest struct {
		Text      string `json:"text" form:"text" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)
