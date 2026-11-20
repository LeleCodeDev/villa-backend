package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	EventBookingRequestQuery struct {
		PaginationQuery
		SortQuery
		Status  model.RequestStatus `form:"status" binding:"omitempty,oneof=pending approved rejected"`
		StartAt string              `form:"start_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		EndAt   string              `form:"end_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		Unpage  bool                `form:"unpage"`
	}

	CreateEventRequest struct {
		Title          string `json:"title" form:"title" binding:"required,max=150"`
		Description    string `json:"description" form:"description" binding:"max=2000"`
		StartAt        string `json:"start_at" form:"start_at" binding:"required,datetime=2006-01-02 15:04"`
		EndAt          string `json:"end_at" form:"end_at" binding:"required,datetime=2006-01-02 15:04"`
		VillaPackageID uint   `json:"villa_package_id" form:"villa_package_id" binding:"required,gt=0"`
	}

	ApproveEventRequest struct {
		ChargeType model.ChargeType `json:"charge_type" binding:"required,oneof=one_time per_hour per_day"`
		Price      int64            `json:"price" binding:"required,gt=0"`
	}

	RejectEventRequest struct {
		Reason string `json:"reason" birding:"required"`
	}

	EventBookingRequestResponse struct {
		ID             uint                  `json:"id"`
		User           UserResponse          `json:"user"`
		VillaPackage   *VillaPackageResponse `json:"villa_package"`
		Title          string                `json:"title"`
		Description    string                `json:"description"`
		StartAt        time.Time             `json:"start_at"`
		EndAt          time.Time             `json:"end_at"`
		Status         model.RequestStatus   `json:"status"`
		ReviewedAt     *time.Time            `json:"reviewed_at"`
		RejectReason   *string               `json:"reject_reason,omitempty"`
		EventBookingID *uint                 `json:"event_booking_id,omitempty"`
		CreatedAt      time.Time             `json:"created_at"`
		UpdatedAt      time.Time             `json:"updated_at"`
	}
)
