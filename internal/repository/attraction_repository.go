package repository

import "gorm.io/gorm"

type AttractionRepository struct {
	db *gorm.DB
}

func NewAttractionRepository(db *gorm.DB) *AttractionRepository {
	return &AttractionRepository{db: db}
}
