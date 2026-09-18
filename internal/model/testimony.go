package model

import (
	"time"

	"gorm.io/gorm"
)

type Testimony struct {
	ID        uint   `gorm:"primaryKey"`
	Star      int    `gorm:"not null"`
	Comment   string `gorm:"not null;type:text"`
	Username  string `gorm:"not null"`
	Order     int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
