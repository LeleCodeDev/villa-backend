package model

import (
	"time"

	"gorm.io/gorm"
)

type VillaPackage struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Subtitle    string `gorm:"not null"`
	MaxCapacity int    `gorm:"not null"`
	Price       int64  `gorm:"not null"`
	SortOrder   int    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
