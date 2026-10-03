package model

import (
	"time"

	"gorm.io/gorm"
)

type Facility struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	LogoID      uint   `gorm:"not null;index"`
	Logo        Logo   `gorm:"constraint:OnDelete:CASCADE"`
	Description string `gorm:"not null"`
	SortOrder   int    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
