package model

import (
	"time"

	"gorm.io/gorm"
)

type EventBookingStatus string

const (
	EventBookingStatusPending   EventBookingStatus = "pending"
	EventBookingStatusPaid      EventBookingStatus = "paid"
	EventBookingStatusCanceled  EventBookingStatus = "canceled"
	EventBookingStatusCompleted EventBookingStatus = "completed"
	EventBookingStatusExpired   EventBookingStatus = "expired"
	// EventBookingStatusRefunded  EventBookingStatus = "refunded"
	EventBookingStatusNotShow EventBookingStatus = "not_show"
)

type EventBooking struct {
	ID             uint               `gorm:"primaryKey"`
	UserID         uint               `gorm:"not null;index"`
	User           User               `gorm:"constraint:OnDelete:CASCADE"`
	VillaPackageID *uint              `gorm:"index"`
	VillaPackage   *VillaPackage      `gorm:"constraint:OnDelete:SET NULL"`
	Title          string             `gorm:"not null"`
	Price          int64              `gorm:"not null"`
	GuestCount     int                `gorm:"not null"`
	Status         EventBookingStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','paid','canceled','completed','expired','not_show')"`
	StartAt        time.Time          `gorm:"not null"`
	EndAt          time.Time          `gorm:"not null"`
	Quantity       int                `gorm:"not null;default:1"`
	TotalPrice     int64              `gorm:"not null"`
	CheckInAt      *time.Time
	CheckOutAt     *time.Time
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
