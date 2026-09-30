package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type SpecificationRepository struct {
	db *gorm.DB
}

func NewSpecificationRepository(db *gorm.DB) *SpecificationRepository {
	return &SpecificationRepository{db: db}
}

func (r *SpecificationRepository) WithTx(tx *gorm.DB) *SpecificationRepository {
	return &SpecificationRepository{db: tx}
}

func (r *SpecificationRepository) GetAll(ctx context.Context, query dto.SpecificationQuery) ([]model.Specification, int64, error) {
	var specifications []model.Specification
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Specification{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&specifications).Error; err != nil {
		return nil, 0, err
	}

	return specifications, total, nil
}

func (r *SpecificationRepository) GetByID(ctx context.Context, id uint) (*model.Specification, error) {
	var specification model.Specification

	if err := r.db.WithContext(ctx).First(&specification, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &specification, nil
}

func (r *SpecificationRepository) Create(ctx context.Context, specification *model.Specification) error {
	return r.db.WithContext(ctx).Create(specification).Error
}

func (r *SpecificationRepository) Update(ctx context.Context, specification *model.Specification) error {
	return r.db.WithContext(ctx).Save(specification).Error
}

func (r *SpecificationRepository) Delete(ctx context.Context, specification *model.Specification) error {
	return r.db.WithContext(ctx).Delete(specification).Error
}

func (r *SpecificationRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Specification{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *SpecificationRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Specification{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *SpecificationRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Specification{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
