package model

import (
	"time"

	"gorm.io/gorm"
)

type RequestStatus string

const (
	RequestPending  RequestStatus = "pending"
	RequestApproved RequestStatus = "approved"
	RequestRejected RequestStatus = "rejected"
)

type EventBookingRequest struct {
	ID             uint          `gorm:"primaryKey"`
	UserID         uint          `gorm:"not null;index"`
	User           User          `gorm:"constraint:OnDelete:CASCADE"`
	Title          string        `gorm:"not null"`
	VillaPackageID *uint         `gorm:"not null;index"`
	VillaPackage   *VillaPackage `gorm:"constraint:OnDelete:SET NULL"`
	Description    string        `gorm:"type:text"`
	StartAt        time.Time     `gorm:"not null"`
	EndAt          time.Time     `gorm:"not null"`
	Status         RequestStatus `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending','approved','rejected')"`
	ReviewedAt     *time.Time
	RejectReason   *string
	EventBookingID *uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
