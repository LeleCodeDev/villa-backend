package model

import (
	"time"

	"gorm.io/gorm"
)

type Pricelist struct {
	ID               uint         `gorm:"primaryKey"`
	RentalOptionID   uint         `gorm:"not null;index"`
	RentalOption     RentalOption `gorm:"constraint:OnDelete:CASCADE"`
	WeekdayPrice     float64      `gorm:"type:numeric(10,2);not null"`
	WeekendPrice     float64      `gorm:"type:numeric(10,2);not null"`
	LongWeekendPrice float64      `gorm:"type:numeric(10,2);not null"`
	HighSeasonPrice  float64      `gorm:"type:numeric(10,2);not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
