package model

import (
	"time"

	"gorm.io/gorm"
)

type HeroSection struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Subtitle  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
