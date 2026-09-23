package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type VillaPackage struct {
	ID          uint            `gorm:"primaryKey"`
	Title       string          `gorm:"not null"`
	Subtitle    string          `gorm:"not null"`
	MaxCapacity int             `gorm:"not null"`
	Price       decimal.Decimal `gorm:"type:numeric(10,2)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
