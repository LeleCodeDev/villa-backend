package model

import (
	"time"

	"gorm.io/gorm"
)

type Faq struct {
	ID        uint   `gorm:"primaryKey"`
	Question  string `gorm:"not null"`
	Answer    string `gorm:"not null"`
	SortOrder int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
