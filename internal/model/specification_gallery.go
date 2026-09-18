package model

import (
	"time"

	"gorm.io/gorm"
)

type SpecificationGallery struct {
	ID        uint    `gorm:"primaryKey"`
	Image1    *string `gorm:"type:varchar(255)"`
	Image2    *string `gorm:"type:varchar(255)"`
	Image3    *string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
