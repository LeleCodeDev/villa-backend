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

func (r *BookingRepository) GetAll(ctx context.Context, query dto.BookingQuery, userID *uint) ([]model.Booking, int64, error) {
	var bookings []model.Booking
	var total int64

	db := r.db.WithContext(ctx).
		Preload("User").
		Preload("RentalOption").
		Preload("Reason").
		Model(&model.Booking{})

	if userID != nil {
		db = db.Where("user_id = ?", userID)
	}

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

func (r *BookingRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.WithContext(ctx).
		Where("payment_order_id = ?", orderID).
		First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) TransitionStatus(ctx context.Context, id uint, from, to model.BookingStatus) (bool, error) {
	updates := map[string]any{"status": to}
	if to == model.BookingStatusPaid {
		updates["paid_at"] = time.Now()
		updates["expires_at"] = nil
	}
	res := r.db.WithContext(ctx).Model(&model.Booking{}).
		Where("id = ? AND status = ?", id, from).
		Updates(updates)
	return res.RowsAffected == 1, res.Error
}

func (r *BookingRepository) SetNullRentalOptionID(ctx context.Context, rentalOptionID uint) error {
	updates := map[string]any{"rental_option_id": nil}
	return r.db.WithContext(ctx).
		Model(&model.Booking{}).
		Where("rental_option_id = ?", rentalOptionID).
		Updates(updates).Error
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
			model.BookingStatusCanceled,
			model.BookingStatusExpired,
			model.BookingStatusCompleted,
			model.BookingStatusNotShow,
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

// func (r *BookingRepository) SetPayment(ctx context.Context, id uint, orderID, token string) error {
// 	return r.db.WithContext(ctx).Model(&model.Booking{}).
// 		Where("id = ?", id).
// 		Updates(map[string]any{"payment_order_id": orderID, "snap_token": token}).Error
// }
