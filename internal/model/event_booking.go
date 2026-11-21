package model

import (
	"time"

	"gorm.io/gorm"
)

type EventBooking struct {
	ID             uint          `gorm:"primaryKey"`
	PaymentOrderID *string       `gorm:"type:varchar(64);uniqueIndex"`
	UserID         uint          `gorm:"not null;index"`
	User           User          `gorm:"constraint:OnDelete:CASCADE"`
	VillaPackageID *uint         `gorm:"index"`
	VillaPackage   *VillaPackage `gorm:"constraint:OnDelete:SET NULL"`
	Status         BookingStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','paid','canceled','completed','expired','not_show')"`
	StartAt        time.Time     `gorm:"not null"`
	EndAt          time.Time     `gorm:"not null"`
	Quantity       int           `gorm:"not null;default:1"`
	TotalPrice     int64         `gorm:"not null"`
	SnapToken      *string       `gorm:"type:varchar(100)"`
	RedirectURL    *string
	PaidAt         *time.Time
	CheckInAt      *time.Time
	CheckOutAt     *time.Time
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
