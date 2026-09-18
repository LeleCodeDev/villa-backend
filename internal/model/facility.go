package model

import (
	"time"

	"gorm.io/gorm"
)

type Facility struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Logo        *string `gorm:"type:varchar(255)"`
	Description string  `gorm:"not null"`
	Order       int     `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
