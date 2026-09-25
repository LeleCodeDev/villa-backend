package model

import (
	"time"

	"gorm.io/gorm"
)

type VillaPackage struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Subtitle    string  `gorm:"not null"`
	MaxCapacity int     `gorm:"not null"`
	Price       float64 `gorm:"type:numeric(10,2)"`
	SortOrder   int     `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
