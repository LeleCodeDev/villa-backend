package model

import (
	"time"

	"gorm.io/gorm"
)

type Specification struct {
	ID        uint   `gorm:"primaryKey"`
	Logo      string `gorm:"not null"`
	Text      string `gorm:"not null"`
	Order     int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
