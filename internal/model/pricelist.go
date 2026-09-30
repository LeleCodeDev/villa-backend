package model

import (
	"time"

	"gorm.io/gorm"
)

type Pricelist struct {
	ID               uint    `gorm:"primaryKey"`
	Title            string  `gorm:"not null"`
	Description      string  `gorm:"not null"`
	MinCapacity      int     `gorm:"not null"`
	MaxCapacity      int     `gorm:"not null"`
	WeekdayPrice     float64 `gorm:"type:numeric(10,2);not null"`
	WeekendPrice     float64 `gorm:"type:numeric(10,2);not null"`
	LongWeekendPrice float64 `gorm:"type:numeric(10,2);not null"`
	HighSeasonPrice  float64 `gorm:"type:numeric(10,2);not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
