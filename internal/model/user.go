package model

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Username  string   `gorm:"not null;index"`
	Email     string   `gorm:"not null;uniqueIndex:idx_email_deleted_at;index"`
	Password  string   `gorm:"not null"`
	Phone     string   `gorm:"not null"`
	Role      UserRole `gorm:"type:varchar(20);not null;default:'user';check:role IN ('admin','user')"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index;uniqueIndex:idx_email_deleted_at"`
}
