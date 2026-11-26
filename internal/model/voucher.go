package model

import (
	"time"

	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypePercent DiscountType = "percent"
	DiscountTypeFixed   DiscountType = "fixed"
)

type Voucher struct {
	ID               uint         `gorm:"primaryKey"`
	Code             string       `gorm:"type:varchar(50);not null;uniqueIndex:idx_code_deleted_at;index"`
	UserID           *uint        `gorm:"index"`
	User             *User        `gorm:"constraint:OnDelete:SET NULL"`
	DiscountType     DiscountType `gorm:"type:varchar(10);not null;check:discount_type IN ('percent','fixed')"`
	DiscountValue    int64        `gorm:"not null"`
	MaxDiscount      *int64
	MinBookingAmount *int64
	UsageLimit       *int
	UsedCount        int `gorm:"not null;default:0"`
	PerUserLimit     *int
	StartAt          *time.Time
	ExpiresAt        *time.Time
	IsActive         bool `gorm:"not null;default:true"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index;uniqueIndex:idx_code_deleted_at;"`
}
