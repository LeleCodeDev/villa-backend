package model

import (
	"time"

	"gorm.io/gorm"
)

type VillaPackageList struct {
	ID             uint         `gorm:"primaryKey"`
	Text           string       `gorm:"not null"`
	VillaPackageID uint         `gorm:"not null;index"`
	VillaPackage   VillaPackage `gorm:"constraint:OnDelete:CASCADE"`
	SortOrder      int          `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
