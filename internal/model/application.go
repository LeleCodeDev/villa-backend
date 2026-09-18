package model

import (
	"time"

	"gorm.io/gorm"
)

type Application struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Logo        *string `gorm:"type:varchar(255)"`
	PhoneNumber string  `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
