package dto

import "time"

type (
	VillaRuleQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	VillaRuleRequest struct {
		Text      string `json:"text" form:"text" binding:"required"`
		SortOrder *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}

	VillaRuleResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (vrq *VillaRuleQuery) SetDefault() {
	vrq.setDefaultPagination()
	vrq.SetDefaultSort(SortAsc)
}
