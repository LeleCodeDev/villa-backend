package dto

import "time"

type (
	VillaPackageListQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	VillaPackageListResponse struct {
		ID        uint      `json:"id"`
		Text      string    `json:"text"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	VillaPackageListRequest struct {
		VillaPackageID uint   `json:"villa_package_id" binding:"required,gt=0"`
		Text           string `json:"text" form:"text" binding:"required"`
		SortOrder      *int   `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (vplq *VillaPackageListQuery) SetDefault() {
	vplq.setDefaultPagination()
	vplq.SetDefaultSort(SortAsc)
}
