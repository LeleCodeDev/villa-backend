package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type TestimonyRepository struct {
	db *gorm.DB
}

func NewTestimonyRepository(db *gorm.DB) *TestimonyRepository {
	return &TestimonyRepository{db: db}
}

func (r *TestimonyRepository) WithTx(tx *gorm.DB) *TestimonyRepository {
	return &TestimonyRepository{db: tx}
}

func (r *TestimonyRepository) GetAll(ctx context.Context, query dto.TestimonyQuery) ([]model.Testimony, int64, error) {
	var testimonies []model.Testimony
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Testimony{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&testimonies).Error; err != nil {
		return nil, 0, err
	}

	return testimonies, total, nil
}

func (r *TestimonyRepository) GetByID(ctx context.Context, id uint) (*model.Testimony, error) {
	var testimony model.Testimony

	if err := r.db.WithContext(ctx).First(&testimony, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &testimony, nil
}

func (r *TestimonyRepository) Create(ctx context.Context, testimony *model.Testimony) error {
	return r.db.WithContext(ctx).Create(testimony).Error
}

func (r *TestimonyRepository) Update(ctx context.Context, testimony *model.Testimony) error {
	return r.db.WithContext(ctx).Save(testimony).Error
}

func (r *TestimonyRepository) Delete(ctx context.Context, testimony *model.Testimony) error {
	return r.db.WithContext(ctx).Delete(testimony).Error
}

func (r *TestimonyRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Testimony{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *TestimonyRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Testimony{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *TestimonyRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Testimony{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
