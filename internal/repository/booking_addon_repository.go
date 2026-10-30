package repository

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type BookingAddonRepository struct {
	db *gorm.DB
}

func NewBookingAddonRepository(db *gorm.DB) *BookingAddonRepository {
	return &BookingAddonRepository{db: db}
}

func (r *BookingAddonRepository) WithTx(tx *gorm.DB) *BookingAddonRepository {
	return &BookingAddonRepository{db: tx}
}

func (r *BookingAddonRepository) GetAllByBookingID(ctx context.Context, bookingID uint) ([]model.BookingAddon, error) {
	var bookingAddons []model.BookingAddon

	if err := r.db.WithContext(ctx).
		Preload("Addon").
		Where("booking_id = ?", bookingID).
		Find(&bookingAddons).Error; err != nil {
		return nil, err
	}

	return bookingAddons, nil
}

func (r *BookingAddonRepository) GetAllInBookingIds(ctx context.Context, bookingIds []uint) ([]model.BookingAddon, error) {
	var bookingAddons []model.BookingAddon

	if err := r.db.WithContext(ctx).
		Preload("Addon").
		Where("booking_id IN ?", bookingIds).
		Find(&bookingAddons).Error; err != nil {
		return nil, err
	}

	return bookingAddons, nil
}
