package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type AddonRepository struct {
	db *gorm.DB
}

func NewAddonRepository(db *gorm.DB) *AddonRepository {
	return &AddonRepository{db: db}
}

func (r *AddonRepository) WithTx(tx *gorm.DB) *AddonRepository {
	return &AddonRepository{db: tx}
}

func (r *AddonRepository) GetAll(ctx context.Context, query dto.AddonQuery) ([]model.Addon, int64, error) {
	var addons []model.Addon
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Addon{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&addons).Error; err != nil {
		return nil, 0, err
	}

	return addons, total, nil
}

func (r *AddonRepository) GetByID(ctx context.Context, id uint) (*model.Addon, error) {
	var addon model.Addon

	if err := r.db.WithContext(ctx).First(&addon, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &addon, nil
}

func (r *AddonRepository) ExistByID(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Addon{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *AddonRepository) Create(ctx context.Context, addon *model.Addon) error {
	return r.db.WithContext(ctx).Create(addon).Error
}

func (r *AddonRepository) Update(ctx context.Context, addon *model.Addon) error {
	return r.db.WithContext(ctx).Save(addon).Error
}

func (r *AddonRepository) Delete(ctx context.Context, addon *model.Addon) error {
	return r.db.WithContext(ctx).Delete(addon).Error
}

func (r *AddonRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Addon{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *AddonRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Addon{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *AddonRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Addon{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
