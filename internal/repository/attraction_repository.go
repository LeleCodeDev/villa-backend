package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type AttractionRepository struct {
	db *gorm.DB
}

func NewAttractionRepository(db *gorm.DB) *AttractionRepository {
	return &AttractionRepository{db: db}
}

func (r *AttractionRepository) WithTx(tx *gorm.DB) *AttractionRepository {
	return &AttractionRepository{db: tx}
}

func (r *AttractionRepository) GetAll(ctx context.Context, query dto.AttractionQuery) ([]model.Attraction, int64, error) {
	var attractions []model.Attraction
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Attraction{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&attractions).Error; err != nil {
		return nil, 0, err
	}

	return attractions, total, nil
}

func (r *AttractionRepository) GetByID(ctx context.Context, id uint) (*model.Attraction, error) {
	var attraction model.Attraction

	if err := r.db.WithContext(ctx).First(&attraction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &attraction, nil
}

func (r *AttractionRepository) ExistByID(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Attraction{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *AttractionRepository) Create(ctx context.Context, attraction *model.Attraction) error {
	return r.db.WithContext(ctx).Create(attraction).Error
}

func (r *AttractionRepository) Update(ctx context.Context, attraction *model.Attraction) error {
	return r.db.WithContext(ctx).Save(attraction).Error
}

func (r *AttractionRepository) Delete(ctx context.Context, attraction *model.Attraction) error {
	return r.db.WithContext(ctx).Delete(attraction).Error
}

func (r *AttractionRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Attraction{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *AttractionRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Attraction{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *AttractionRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Attraction{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
