package model

import (
	"time"

	"gorm.io/gorm"
)

type GalleryType string

const (
	GalleryTypeBooking GalleryType = "booking"
	GalleryTypeWedding GalleryType = "wedding"
)

type Gallery struct {
	ID          uint        `gorm:"primaryKey"`
	Title       string      `gorm:"not null"`
	Description string      `gorm:"not null"`
	Image       *string     `gorm:"type:varchar(255)"`
	Type        GalleryType `gorm:"type:varchar(20);not null;default:'booking';check:type IN ('booking','wedding')"`
	SortOrder   int         `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
