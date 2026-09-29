package dto

type SortType string

const (
	SortDesc SortType = "desc"
	SortAsc  SortType = "asc"
)

type SortQuery struct {
	Sort SortType `json:"sort" form:"sort" binding:"omitempty,oneof=asc desc"`
}

func (sq *SortQuery) SetDefaultSort(sort SortType) {
	if sq.Sort == "" {
		sq.Sort = sort
	}
}
