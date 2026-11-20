package model

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending       BookingStatus = "pending"
	BookingStatusPaid          BookingStatus = "paid"
	BookingStatusPartiallyPaid BookingStatus = "partially_paid"
	BookingStatusCanceled      BookingStatus = "canceled"
	BookingStatusCompleted     BookingStatus = "completed"
	BookingStatusExpired       BookingStatus = "expired"
	// BookingStatusRefunded  BookingStatus = "refunded"
	BookingStatusNotShow BookingStatus = "not_show"
)

type PaymentType string

const (
	PaymentTypeFull    PaymentType = "full"
	PaymentTypeDeposit PaymentType = "deposit"
)

type Booking struct {
	ID             uint          `gorm:"primaryKey"`
	UserID         uint          `gorm:"not null;index"`
	User           User          `gorm:"constraint:OnDelete:CASCADE"`
	GuestCount     int           `gorm:"not null"`
	RentalOptionID *uint         `gorm:"index"`
	RentalOption   *RentalOption `gorm:"constraint:OnDelete:SET NULL"`
	ReasonID       uint          `gorm:"not null;index"`
	Reason         Reason        `gorm:"constraint:OnDelete:CASCADE"`
	Status         BookingStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'partially_paid', 'paid','canceled','completed','expired','not_show')"`
	PaymentType    PaymentType   `gorm:"type:varchar(20);not null;default:'full'"`
	CheckInDate    time.Time     `gorm:"not null;type:date"`
	CheckOutDate   time.Time     `gorm:"not null;type:date"`
	BookingPrice   int64         `gorm:"not null"`
	TotalPrice     int64         `gorm:"not null"`
	DepositAmount  *int64
	BalanceDueAt   *time.Time
	CheckInAt      *time.Time
	CheckOutAt     *time.Time
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
