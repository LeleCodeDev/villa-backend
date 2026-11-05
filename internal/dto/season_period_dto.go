package dto

import "time"

type (
	SeasonPeriodQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	SeasonPeriodRequest struct {
		Name      string `json:"name" form:"name" binding:"required"`
		StartDate string `json:"start_date" form:"start_date" binding:"required,datetime=2006-01-02"`
		EndDate   string `json:"end_date" form:"end_date" binding:"required,datetime=2006-01-02"`
	}

	SeasonPeriodResponse struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (spq *SeasonPeriodQuery) SetDefault() {
	spq.setDefaultPagination()
	spq.SetDefaultSort(SortAsc)
}
