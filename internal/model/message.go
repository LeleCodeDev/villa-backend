package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	Firstname string `gorm:"not null;type:text"`
	Lastname  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Message   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
