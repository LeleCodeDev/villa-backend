package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) WithTx(tx *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: tx}
}

func (r *PaymentRepository) GetAllByBookingID(ctx context.Context, bookingID uint) ([]model.Payment, error) {
	var payments []model.Payment

	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetAllByBookingIDAndStatus(ctx context.Context, bookingID uint, status model.PaymentStatus) ([]model.Payment, error) {
	var payments []model.Payment

	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Where("status = ?", status).
		Order("created_at ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetFirstByBookingIDAndStatusAndKind(ctx context.Context, bookingID uint, status model.PaymentStatus, kind model.PaymentKind) (*model.Payment, error) {
	var payment model.Payment

	if err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Where("kind = ?", kind).
		Where("status = ?", status).
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) GetAllInBookingIDs(ctx context.Context, bookingIds []uint) ([]model.Payment, error) {
	var payments []model.Payment

	if err := r.db.WithContext(ctx).
		Where("booking_id IN ?", bookingIds).
		Order("created_at ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetAllByEventBookingID(ctx context.Context, eventBookingID uint) ([]model.Payment, error) {
	var payments []model.Payment

	if err := r.db.WithContext(ctx).
		Where("event_booking_id = ?", eventBookingID).
		Order("created_at ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetByEventBookingID(ctx context.Context, eventBookingID uint) (*model.Payment, error) {
	var payment model.Payment

	if err := r.db.WithContext(ctx).
		Where("event_booking_id = ?", eventBookingID).
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) GetByEventBookingIDAndStatus(ctx context.Context, eventBookingID uint, status model.PaymentStatus) (*model.Payment, error) {
	var payment model.Payment

	if err := r.db.WithContext(ctx).
		Where("event_booking_id = ?", eventBookingID).
		Where("status = ?", status).
		Order("created_at DESC").
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) GetAllInEventBookingIds(ctx context.Context, eventBookingIds []uint) ([]model.Payment, error) {
	var payments []model.Payment

	if err := r.db.WithContext(ctx).
		Where("event_booking_id IN ?", eventBookingIds).
		Order("created_at ASC").
		Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).
		Preload("Booking").
		Preload("Booking.User").
		Preload("Booking.RentalOption").
		Preload("Booking.BookingAddons.Addon").
		Preload("EventBooking").
		Preload("EventBooking.User").
		Preload("EventBooking.VillaPackage").
		Where("order_id = ?", orderID).
		First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *PaymentRepository) Delete(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Delete(payment).Error
}

func (r *PaymentRepository) TransitionStatus(ctx context.Context, id uint, from, to model.PaymentStatus) (bool, error) {
	updates := map[string]any{"status": to}
	if to == model.PaymentStatusPaid {
		updates["paid_at"] = time.Now()
	}

	res := r.db.WithContext(ctx).Model(&model.Payment{}).
		Where("id = ? AND status = ?", id, from).
		Updates(updates)

	return res.RowsAffected == 1, res.Error
}
