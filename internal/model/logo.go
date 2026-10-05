package model

import (
	"time"

	"gorm.io/gorm"
)

type Logo struct {
	ID        uint   `gorm:"primaryKey"`
	Image     string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
