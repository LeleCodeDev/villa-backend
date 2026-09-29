package dto

type PaginationQuery struct {
	Page int `json:"page" form:"page" binding:"omitempty,gt=0"`
	Size int `json:"size" form:"size" binding:"omitempty,gt=0"`
}

func (pq *PaginationQuery) setDefaultPagination() {
	if pq.Page < 1 {
		pq.Page = 1
	}

	if pq.Size < 1 {
		pq.Size = 10
	}
}

func (pq *PaginationQuery) GetOffset() int {
	return (pq.Page - 1) * pq.Size
}
