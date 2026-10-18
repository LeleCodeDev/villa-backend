package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
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

func (r *PricelistRepository) GetAll(ctx context.Context, query dto.PricelistQuery) ([]model.Pricelist, int64, error) {
	var pricelists []model.Pricelist
	var total int64

	db := r.db.WithContext(ctx).Model(model.Pricelist{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at" + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&pricelists).Error; err != nil {
		return nil, 0, err
	}

	return pricelists, total, nil
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
