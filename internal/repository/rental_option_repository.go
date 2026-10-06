package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type RentalOptionRepository struct {
	db *gorm.DB
}

func NewRentalOptionRepository(db *gorm.DB) *RentalOptionRepository {
	return &RentalOptionRepository{db: db}
}

func (r *RentalOptionRepository) WithTx(tx *gorm.DB) *RentalOptionRepository {
	return &RentalOptionRepository{db: tx}
}

func (r *RentalOptionRepository) GetAll(ctx context.Context, query dto.RentalOptionQuery) ([]model.RentalOption, int64, error) {
	var rentalOptions []model.RentalOption
	var total int64

	db := r.db.WithContext(ctx).Model(&model.RentalOption{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&rentalOptions).Error; err != nil {
		return nil, 0, err
	}

	return rentalOptions, total, nil
}

func (r *RentalOptionRepository) GetByID(ctx context.Context, id uint) (*model.RentalOption, error) {
	var rentalOption model.RentalOption

	if err := r.db.WithContext(ctx).First(&rentalOption, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &rentalOption, nil
}

func (r *RentalOptionRepository) Create(ctx context.Context, rentalOption *model.RentalOption) error {
	return r.db.WithContext(ctx).Create(rentalOption).Error
}

func (r *RentalOptionRepository) Update(ctx context.Context, rentalOption *model.RentalOption) error {
	return r.db.WithContext(ctx).Save(rentalOption).Error
}

func (r *RentalOptionRepository) Delete(ctx context.Context, rentalOption *model.RentalOption) error {
	return r.db.WithContext(ctx).Delete(rentalOption).Error
}

func (r *RentalOptionRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.RentalOption{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *RentalOptionRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.RentalOption{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *RentalOptionRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.RentalOption{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
