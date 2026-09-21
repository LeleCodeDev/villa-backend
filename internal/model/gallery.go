package model

import (
	"time"

	"gorm.io/gorm"
)

type Gallery struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Image       *string `gorm:"type:varchar(255)"`
	SortOrder   int     `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
