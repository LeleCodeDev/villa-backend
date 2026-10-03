package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type FacilityRepository struct {
	db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) *FacilityRepository {
	return &FacilityRepository{db: db}
}

func (r *FacilityRepository) WithTx(tx *gorm.DB) *FacilityRepository {
	return &FacilityRepository{db: tx}
}

func (r *FacilityRepository) GetAll(ctx context.Context, query dto.FacilityQuery) ([]model.Facility, int64, error) {
	var facilities []model.Facility
	var total int64

	db := r.db.WithContext(ctx).Preload("Logo").Model(&model.Facility{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&facilities).Error; err != nil {
		return nil, 0, err
	}

	return facilities, total, nil
}

func (r *FacilityRepository) GetByID(ctx context.Context, id uint) (*model.Facility, error) {
	var facility model.Facility

	if err := r.db.WithContext(ctx).Preload("Logo").First(&facility, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &facility, nil
}

func (r *FacilityRepository) Create(ctx context.Context, facility *model.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

func (r *FacilityRepository) Update(ctx context.Context, facility *model.Facility) error {
	return r.db.WithContext(ctx).Save(facility).Error
}

func (r *FacilityRepository) Delete(ctx context.Context, facility *model.Facility) error {
	return r.db.WithContext(ctx).Delete(facility).Error
}

func (r *FacilityRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Facility{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *FacilityRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Facility{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *FacilityRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Facility{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
