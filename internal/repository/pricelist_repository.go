package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type PricelistRepository struct {
	db *gorm.DB
}

func NewPricelistRepository(db *gorm.DB) *PricelistRepository {
	return &PricelistRepository{db: db}
}

func (r *PricelistRepository) WithTx(tx *gorm.DB) *PricelistRepository {
	return &PricelistRepository{db: tx}
}

func (r *PricelistRepository) GetAll(ctx context.Context) ([]model.Pricelist, error) {
	var pricelists []model.Pricelist

	if err := r.db.WithContext(ctx).Find(&pricelists).Error; err != nil {
		return nil, err
	}

	return pricelists, nil
}

func (r *PricelistRepository) GetByID(ctx context.Context, id uint) (*model.Pricelist, error) {
	var pricelist model.Pricelist

	if err := r.db.WithContext(ctx).First(&pricelist, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &pricelist, nil
}

func (r *PricelistRepository) Create(ctx context.Context, pricelist *model.Pricelist) error {
	return r.db.WithContext(ctx).Create(pricelist).Error
}

func (r *PricelistRepository) Update(ctx context.Context, pricelist *model.Pricelist) error {
	return r.db.WithContext(ctx).Save(pricelist).Error
}

func (r *PricelistRepository) Delete(ctx context.Context, pricelist *model.Pricelist) error {
	return r.db.WithContext(ctx).Delete(pricelist).Error
}

func (r *PricelistRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Pricelist{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *PricelistRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Pricelist{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *PricelistRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Pricelist{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
