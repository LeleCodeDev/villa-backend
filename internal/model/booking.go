package model

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusCanceled  BookingStatus = "canceled"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusExpired   BookingStatus = "expired"
	BookingStatusRefunded  BookingStatus = "refunded"
	BookingStatusNotShow   BookingStatus = "not_show"
)

type Booking struct {
	ID             uint          `gorm:"primaryKey"`
	UserID         uint          `gorm:"not null;index"`
	User           User          `gorm:"constraint:OnDelete:CASCADE"`
	GuestCount     int           `gorm:"not null"`
	CheckIn        time.Time     `gorm:"not null;type:date"`
	CheckOut       time.Time     `gorm:"type:date"`
	RentalOptionID uint          `gorm:"not null;index"`
	RentalOption   RentalOption  `gorm:"constraint:OnDelete:CASCADE"`
	ReasonID       uint          `gorm:"not null;index"`
	Reason         Reason        `gorm:"constraint:OnDelete:CASCADE"`
	Status         BookingStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','paid','canceled','completed','expired','refunded','not_show')"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
