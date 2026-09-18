package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type HeroSectionRepository struct {
	db *gorm.DB
}

func NewHeroSectionRepository(db *gorm.DB) *HeroSectionRepository {
	return &HeroSectionRepository{db: db}
}

func (r *HeroSectionRepository) WithTx(tx *gorm.DB) *HeroSectionRepository {
	return &HeroSectionRepository{db: tx}
}

func (r *HeroSectionRepository) Get(ctx context.Context) (*model.HeroSection, error) {
	var heroSection model.HeroSection

	if err := r.db.WithContext(ctx).Model(&model.HeroSection{}).First(&heroSection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &heroSection, nil
}

func (r *HeroSectionRepository) Exist(ctx context.Context) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.HeroSection{}).
		Count(&count).Error

	return count > 0, err
}

func (r *HeroSectionRepository) Create(ctx context.Context, heroSection *model.HeroSection) error {
	return r.db.WithContext(ctx).Create(heroSection).Error
}

func (r *HeroSectionRepository) Update(ctx context.Context, heroSection *model.HeroSection) error {
	return r.db.WithContext(ctx).Save(heroSection).Error
}
