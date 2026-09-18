package model

import (
	"time"

	"gorm.io/gorm"
)

type AddonList struct {
	ID        uint   `gorm:"primaryKey"`
	Text      string `gorm:"not null"`
	AddonID   uint   `gorm:"not null;index"`
	Addon     Addon  `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
