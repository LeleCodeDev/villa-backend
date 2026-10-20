package dto

import "time"

type (
	ReasonQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

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

func (rq *ReasonQuery) SetDefault() {
	rq.setDefaultPagination()
	rq.SetDefaultSort(SortAsc)
}
