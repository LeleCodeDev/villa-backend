package mapper

import (
	"strings"
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToVoucherResponse(voucher *model.Voucher) dto.VoucherResponse {
	return dto.VoucherResponse{
		ID:               voucher.ID,
		Code:             voucher.Code,
		DiscountType:     voucher.DiscountType,
		DiscountValue:    voucher.DiscountValue,
		MaxDiscount:      voucher.MaxDiscount,
		MinBookingAmount: voucher.MinBookingAmount,
		UsageLimit:       voucher.UsageLimit,
		UsedCount:        voucher.UsedCount,
		PerUserLimit:     voucher.PerUserLimit,
		StartAt:          voucher.StartAt,
		ExpiresAt:        voucher.ExpiresAt,
		IsActive:         voucher.IsActive,

		CreatedAt: voucher.CreatedAt,
		UpdatedAt: voucher.UpdatedAt,
	}
}

func ToVoucherModel(req dto.VoucherRequest, startAt, expiresAt *time.Time) *model.Voucher {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &model.Voucher{
		Code:             normalizeVoucherCode(req.Code),
		DiscountType:     req.DiscountType,
		DiscountValue:    req.DiscountValue,
		MaxDiscount:      req.MaxDiscount,
		MinBookingAmount: req.MinBookingAmount,
		UsageLimit:       req.UsageLimit,
		PerUserLimit:     req.PerUserLimit,
		StartAt:          startAt,
		ExpiresAt:        expiresAt,
		IsActive:         isActive,
		SourceBookingID:  nil,
	}
}

func UpdateVoucherModel(voucher *model.Voucher, req dto.VoucherRequest, startsAt, expiresAt *time.Time) {
	voucher.Code = normalizeVoucherCode(req.Code)
	voucher.DiscountType = req.DiscountType
	voucher.DiscountValue = req.DiscountValue
	voucher.MaxDiscount = req.MaxDiscount
	voucher.MinBookingAmount = req.MinBookingAmount
	voucher.UsageLimit = req.UsageLimit
	voucher.PerUserLimit = req.PerUserLimit
	voucher.StartAt = startsAt
	voucher.ExpiresAt = expiresAt

	if req.IsActive != nil {
		voucher.IsActive = *req.IsActive
	}
}

func normalizeVoucherCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
