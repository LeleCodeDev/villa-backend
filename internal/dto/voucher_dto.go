package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	VoucherQuery struct {
		PaginationQuery
		SortQuery
		Search   string `form:"search"`
		IsActive *bool  `form:"is_active"`
		Unpage   bool   `form:"unpage"`
	}

	VoucherRequest struct {
		Code             string             `json:"code" form:"code" binding:"required,min=3,max=50,alphanum"`
		DiscountType     model.DiscountType `json:"discount_type" form:"discount_type" binding:"required,oneof=percent fixed"`
		DiscountValue    int64              `json:"discount_value" form:"discount_value" binding:"required,gt=0"`
		MaxDiscount      *int64             `json:"max_discount" form:"max_discount" binding:"omitempty,gt=0"`
		MinBookingAmount *int64             `json:"min_booking_amount" form:"min_booking_amount" binding:"omitempty,gt=0"`
		UsageLimit       *int               `json:"usage_limit" form:"usage_limit" binding:"omitempty,gt=0"`
		PerUserLimit     *int               `json:"per_user_limit" form:"per_user_limit" binding:"omitempty,gt=0"`
		StartAt          *string            `json:"start_at" form:"start_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		ExpiresAt        *string            `json:"expires_at" form:"expires_at" binding:"omitempty,datetime=2006-01-02 15:04"`
		IsActive         *bool              `json:"is_active" form:"is_active"`
	}

	VoucherResponse struct {
		ID               uint               `json:"id"`
		Code             string             `json:"code"`
		DiscountType     model.DiscountType `json:"discount_type"`
		DiscountValue    int64              `json:"discount_value"`
		MaxDiscount      *int64             `json:"max_discount,omitempty"`
		MinBookingAmount *int64             `json:"min_booking_amount,omitempty"`
		UsageLimit       *int               `json:"usage_limit,omitempty"`
		UsedCount        int                `json:"used_count"`
		PerUserLimit     *int               `json:"per_user_limit,omitempty"`
		StartAt          *time.Time         `json:"start_at,omitempty"`
		ExpiresAt        *time.Time         `json:"expires_at,omitempty"`
		IsActive         bool               `json:"is_active"`
		CreatedAt        time.Time          `json:"created_at"`
		UpdatedAt        time.Time          `json:"updated_at"`
	}

	VoucherValidateRequest struct {
		Code   string `json:"code" form:"code" binding:"required,min=3,max=50"`
		Amount int64  `json:"amount" form:"amount" binding:"required,gt=0"`
	}

	VoucherValidateResponse struct {
		Code           string `json:"code"`
		DiscountAmount int64  `json:"discount_amount"`
		FinalAmount    int64  `json:"final_amount"`
	}
)

func (vq *VoucherQuery) SetDefault() {
	vq.setDefaultPagination()
	vq.SetDefaultSort(SortDesc)
}
