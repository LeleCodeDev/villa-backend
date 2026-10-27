package model

import (
	"time"

	"gorm.io/gorm"
)

type SeasonPeriod struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	StartDate time.Time `gorm:"not null;type:date;index"`
	EndDate   time.Time `gorm:"not null;type:date"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
