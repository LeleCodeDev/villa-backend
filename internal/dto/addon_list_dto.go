package dto

import "time"

type (
	AddonListQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	AddonListResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	AddonListRequest struct {
		AddonID   uint   `json:"addon_id" binding:"required,gt=0"`
		Text      string `json:"text" form:"text" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (alq *AddonListQuery) SetDefault() {
	alq.setDefaultPagination()
	alq.SetDefaultSort(SortAsc)
}
