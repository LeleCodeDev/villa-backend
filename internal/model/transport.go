package model

import (
	"time"

	"gorm.io/gorm"
)

type Transport struct {
	ID           uint       `gorm:"primaryKey"`
	AttractionID uint       `gorm:"not null;index"`
	Attraction   Attraction `gorm:"constraint:OnDelete:CASCADE"`
	Vehicle      string     `gorm:"not null"`
	Distance     string     `gorm:"not null"`
	TimeRequired string     `gorm:"not null"`
	SortOrder    int        `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
