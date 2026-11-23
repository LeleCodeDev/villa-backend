package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	VillaPackageQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool `form:"unpage"`
	}

	VillaPackageResponse struct {
		ID               uint                       `json:"id"`
		Title            string                     `json:"title"`
		Subtitle         string                     `json:"subtitle"`
		Type             model.ChargeType           `json:"type"`
		Price            int64                      `json:"price"`
		Custom           bool                       `json:"custom"`
		RequiresApproval bool                       `json:"requires_approval"`
		MaxCapacity      int                        `json:"max_capacity"`
		IsWedding        bool                       `json:"is_wedding"`
		Lists            []VillaPackageListResponse `json:"lists,omitempty"`
		SortOrder        int                        `json:"sort_order"`
		CreatedAt        time.Time                  `json:"created_at"`
		UpdatedAt        time.Time                  `json:"updated_at"`
	}

	VillaPackageRequest struct {
		Title            string           `json:"title" form:"title" binding:"required"`
		Subtitle         string           `json:"subtitle" form:"subtitle" binding:"required"`
		Type             model.ChargeType `json:"type" form:"type" binding:"required,oneof=per_hour per_day one_time"`
		MaxCapacity      int              `json:"max_capacity" form:"max_capacity" binding:"required,gt=0"`
		Price            int64            `json:"price" form:"price" binding:"required,gt=0"`
		IsWedding        bool             `json:"is_wedding" form:"is_wedding"`
		RequiresApproval *bool            `json:"requires_approval" form:"requires_approval" binding:"required"`
		SortOrder        *int             `json:"sort_order" form:"sort_order" binding:"omitempty,gt=0"`
	}
)

func (vpq *VillaPackageQuery) SetDefault() {
	vpq.setDefaultPagination()
	vpq.SetDefaultSort(SortAsc)
}
