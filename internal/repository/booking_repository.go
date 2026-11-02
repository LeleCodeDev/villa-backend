package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) WithTx(tx *gorm.DB) *BookingRepository {
	return &BookingRepository{db: tx}
}

func (r *BookingRepository) GetAll(ctx context.Context, query dto.BookingQuery) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64

	db := r.db.WithContext(ctx).
		Preload("User").
		Preload("RentalOption").
		Preload("Reason").
		Model(&model.Booking{})

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.StartDate != "" {
		db = db.Where("check_in_date >= ? ", query.StartDate)
	}

	if query.EndDate != "" {
		db = db.Where("check_out_date <= ? ", query.EndDate)
	}

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *BookingRepository) GetByID(ctx context.Context, id uint) (*model.Booking, error) {
	var booking model.Booking

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("RentalOption").
		Preload("Reason").
		First(&booking, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

func (r *BookingRepository) Update(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}

func (r *BookingRepository) Delete(ctx context.Context, booking *model.Booking) error {
	return r.db.WithContext(ctx).Delete(booking).Error
}

func (r *BookingRepository) GetBookingConflicts(ctx context.Context, optionID uint, in time.Time, out time.Time, excludeID uint) ([]model.Booking, error) {
	var conflicts []model.Booking
	db := r.db.WithContext(ctx).
		Where("rental_option_id = ?", optionID).
		Where("status NOT IN ?", []model.BookingStatus{
			model.BookingStatusCanceled, model.BookingStatusExpired,
		}).
		Where("check_in_date <= ? AND check_out_date >= ?", out, in)

	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}

	if err := db.Find(&conflicts).Error; err != nil {
		return nil, err
	}

	return conflicts, nil
}
