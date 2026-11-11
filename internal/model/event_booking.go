package model

import (
	"time"

	"gorm.io/gorm"
)

type EventBookingType string

const (
	EventBookingPerHour EventBookingType = "per_hour"
	EventBookingPerDay  EventBookingType = "per_day"
	EventBookingOneTime EventBookingType = "one_time"
)

type EventBooking struct {
	ID             uint             `gorm:"primaryKey"`
	PaymentOrderID *string          `gorm:"type:varchar(64);uniqueIndex"`
	UserID         uint             `gorm:"not null;index"`
	User           User             `gorm:"constraint:OnDelete:CASCADE"`
	Status         BookingStatus    `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','paid','canceled','completed','expired','not_show')"`
	Type           EventBookingType `gorm:"type:varchar(20);not null;default:'one_time';check:type IN ('one_time','per_hour','per_day')"`
	StartAt        time.Time        `gorm:"not null"`
	EndAt          time.Time        `gorm:"not null"`
	Price          int64            `gorm:"not null"`
	Quantity       int              `gorm:"not null;default:1"`
	TotalPrice     int64            `gorm:"not null"`
	SnapToken      *string          `gorm:"type:varchar(100)"`
	RedirectURL    *string
	PaidAt         *time.Time
	CheckInAt      *time.Time
	CheckOutAt     *time.Time
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
