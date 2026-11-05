package model

import (
	"time"

	"gorm.io/gorm"
)

type (
	SeasonType string
	DayType    string
)

const (
	SeasonReguler SeasonType = "reguler"
	SeasonHigh    SeasonType = "high"

	DayWeekday DayType = "weekday"
	DayWeekend DayType = "weekend"
)

type RentalPrice struct {
	ID             uint         `gorm:"primaryKey"`
	RentalOptionID uint         `gorm:"not null;uniqueIndex:idx_price,priority:1,where:deleted_at IS NULL"`
	RentalOption   RentalOption `gorm:"constraint:OnDelete:CASCADE"`
	Season         SeasonType   `gorm:"type:varchar(10);not null;uniqueIndex:idx_price,priority:2,where:deleted_at IS NULL"`
	DayType        DayType      `gorm:"type:varchar(10);not null;uniqueIndex:idx_price,priority:3,where:deleted_at IS NULL"`
	Price          int64        `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
