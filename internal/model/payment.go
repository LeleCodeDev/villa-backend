package model

import (
	"time"
)

type PaymentKind string

const (
	PaymentKindFull    PaymentKind = "full"
	PaymentKindDeposit PaymentKind = "deposit"
	PaymentKindBalance PaymentKind = "balance"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusCanceled PaymentStatus = "canceled"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusExpired  PaymentStatus = "expired"
)

type Payment struct {
	ID             uint          `gorm:"primaryKey"`
	OrderID        string        `gorm:"type:varchar(64);not null;uniqueIndex"`
	BookingID      *uint         `gorm:"index"`
	Booking        *Booking      `gorm:"constraint:OnDelete:RESTRICT"`
	EventBookingID *uint         `gorm:"index"`
	EventBooking   *EventBooking `gorm:"constraint:OnDelete:RESTRICT"`
	Kind           PaymentKind   `gorm:"type:varchar(20);not null;default:'full';check:kind IN ('full', 'deposit', 'balance')"`
	Status         PaymentStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'paid', 'failed', 'expired')"`
	Amount         int64         `gorm:"not null"`
	SnapToken      *string       `gorm:"type:varchar(100)"`
	RedirectURL    *string
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
