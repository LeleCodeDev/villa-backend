package model

import (
	"time"

	"gorm.io/gorm"
)

type Specification struct {
	ID        uint   `gorm:"primaryKey"`
	Text      string `gorm:"not null"`
	LogoID    uint   `gorm:"not null;index"`
	Logo      Logo   `gorm:"constraint:OnDelete:CASCADE"`
	SortOrder int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
