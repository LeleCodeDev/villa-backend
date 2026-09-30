package model

import (
	"time"

	"gorm.io/gorm"
)

type Facility struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Logo        string `gorm:"not null"`
	Description string `gorm:"not null"`
	SortOrder   int    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
