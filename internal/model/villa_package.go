package model

import (
	"time"

	"gorm.io/gorm"
)

type ChargeType string

const (
	ChargePerHour ChargeType = "per_hour"
	ChargePerDay  ChargeType = "per_day"
	ChargeOneTime ChargeType = "one_time"
)

type VillaPackage struct {
	ID        uint       `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	Subtitle  string     `gorm:"not null"`
	Price     int64      `gorm:"not null"`
	Type      ChargeType `gorm:"type:varchar(20);not null;default:'one_time';check:type IN ('one_time','per_hour','per_day')"`
	Custom    bool       `gorm:"not null"`
	SortOrder int        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
