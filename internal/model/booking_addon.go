package model

import (
	"time"

	"gorm.io/gorm"
)

type BookingAddon struct {
	ID        uint    `gorm:"primaryKey"`
	BookingID uint    `gorm:"not null;index"`
	Booking   Booking `gorm:"constraint:OnDelete:CASCADE"`
	AddonID   uint    `gorm:"not null;index"`
	Addon     Addon   `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
