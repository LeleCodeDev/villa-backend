package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type LogoRepository struct {
	db *gorm.DB
}

func NewLogoRepository(db *gorm.DB) *LogoRepository {
	return &LogoRepository{db: db}
}

func (r *LogoRepository) WithTx(tx *gorm.DB) *LogoRepository {
	return &LogoRepository{db: tx}
}

func (r *LogoRepository) GetAll(ctx context.Context, query dto.LogoQuery) ([]model.Logo, int64, error) {
	var logos []model.Logo
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Logo{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&logos).Error; err != nil {
		return nil, 0, err
	}

	return logos, total, nil
}

func (r *LogoRepository) GetByID(ctx context.Context, id uint) (*model.Logo, error) {
	var logo model.Logo

	if err := r.db.WithContext(ctx).First(&logo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &logo, nil
}

func (r *LogoRepository) Create(ctx context.Context, logo *model.Logo) error {
	return r.db.WithContext(ctx).Create(logo).Error
}

func (r *LogoRepository) Update(ctx context.Context, logo *model.Logo) error {
	return r.db.WithContext(ctx).Save(logo).Error
}

func (r *LogoRepository) Delete(ctx context.Context, logo *model.Logo) error {
	return r.db.WithContext(ctx).Delete(logo).Error
}
