package dto

import "time"

type (
	LogoQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	LogoResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	LogoRequest struct {
		Text string `json:"text" form:"text" binding:"required"`
	}
)

func (lq *LogoQuery) SetDefault() {
	lq.setDefaultPagination()
	lq.SetDefaultSort(SortAsc)
}
