package model

import (
	"time"

	"gorm.io/gorm"
)

type Reason struct {
	ID        uint   `gorm:"primaryKey"`
	Text      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
