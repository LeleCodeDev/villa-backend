package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Pricelist struct {
	ID               uint            `gorm:"primaryKey"`
	Title            string          `gorm:"not null"`
	Description      string          `gorm:"not null"`
	MinCapacity      int             `gorm:"not null"`
	MaxCapacity      int             `gorm:"not null"`
	WeekdayPrice     decimal.Decimal `gorm:"type:numeric(10,2)"`
	WeekendPrice     decimal.Decimal `gorm:"type:numeric(10,2)"`
	LongWeekendPrice decimal.Decimal `gorm:"type:numeric(10,2)"`
	HighSeasonPrice  decimal.Decimal `gorm:"type:numeric(10,2)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
