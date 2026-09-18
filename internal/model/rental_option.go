package model

import (
	"time"

	"gorm.io/gorm"
)

type RentalOption struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Subtitle    string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	RoomCount   int     `gorm:"not null"`
	MaxCapacity int     `gorm:"not null"`
	Order       int     `gorm:"not null"`
	Image       *string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
