package repository

import (
	"context"
	"errors"
	"time"

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

func (r *RentalOptionRepository) GetAvailable(ctx context.Context, checkIn time.Time, checkOut time.Time, guests int) ([]model.RentalOption, error) {
	var options []model.RentalOption

	conflicting := r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Select("rental_option_id").
		Where("rental_option_id IS NOT NULL").
		Where("check_in_date <= ? AND check_out_date >= ?", checkOut, checkIn).
		Where(
			"status = ? OR (status = ? AND (expires_at IS NULL OR expires_at > ?))",
			model.BookingStatusPaid,
			model.BookingStatusPending,
			time.Now(),
		)

	err := r.db.WithContext(ctx).
		Where("max_capacity >= ?", guests).
		Where("id NOT IN (?)", conflicting).
		Order("sort_order ASC").
		Find(&options).Error

	return options, err
}
