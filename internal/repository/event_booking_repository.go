package repository

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type EventBookingRepository struct {
	db *gorm.DB
}

func NewEventBookingRepository(db *gorm.DB) *EventBookingRepository {
	return &EventBookingRepository{db: db}
}

func (r *EventBookingRepository) WithTx(tx *gorm.DB) *EventBookingRepository {
	return &EventBookingRepository{db: tx}
}

func (r *EventBookingRepository) GetAll(ctx context.Context, query dto.EventBookingQuery, userID *uint) ([]model.EventBooking, int64, error) {
	var eventBookings []model.EventBooking
	var total int64

	db := r.db.WithContext(ctx).
		Preload("User").
		Preload("VillaPackage").
		Model(&model.EventBooking{})

	if query.Type != "" {
		db = db.
			Joins("JOIN villa_packages ON villa_packages.id = event_bookings.villa_package_id AND villa_packages.deleted_at IS NULL").
			Where("villa_packages.type = ? ", query.Type)
	}

	if userID != nil {
		db = db.Where("event_bookings.user_id = ?", userID)
	}

	if query.Status != "" {
		db = db.Where("event_bookings.status = ?", query.Status)
	}

	if query.StartAt != "" {
		start, err := time.ParseInLocation("2006-01-02 15:04", query.StartAt, config.Env.Location)
		if err != nil {
			return nil, 0, err
		}

		db = db.Where("event_bookings.start_at >= ?", start)
	}

	if query.EndAt != "" {
		end, err := time.ParseInLocation("2006-01-02 15:04", query.EndAt, config.Env.Location)
		if err != nil {
			return nil, 0, err
		}

		db = db.Where("event_bookings.end_at <= ?", end)
	}

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("event_bookings.created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&eventBookings).Error; err != nil {
		return nil, 0, err
	}

	return eventBookings, total, nil
}

func (r *EventBookingRepository) GetByID(ctx context.Context, id uint) (*model.EventBooking, error) {
	var eventBookings model.EventBooking

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("VillaPackage").
		First(&eventBookings, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &eventBookings, nil
}

func (r *EventBookingRepository) Create(ctx context.Context, eventBooking *model.EventBooking) error {
	return r.db.WithContext(ctx).Create(eventBooking).Error
}

func (r *EventBookingRepository) Update(ctx context.Context, eventBooking *model.EventBooking) error {
	return r.db.WithContext(ctx).Save(eventBooking).Error
}

func (r *EventBookingRepository) Delete(ctx context.Context, eventBooking *model.EventBooking) error {
	return r.db.WithContext(ctx).Delete(eventBooking).Error
}

func (r *EventBookingRepository) GetByOrderID(ctx context.Context, orderID string) (*model.EventBooking, error) {
	var eventBooking model.EventBooking
	if err := r.db.WithContext(ctx).
		Where("payment_order_id = ?", orderID).
		First(&eventBooking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return &eventBooking, nil
}

func (r *EventBookingRepository) TransitionStatus(ctx context.Context, id uint, from, to model.BookingStatus) (bool, error) {
	updates := map[string]any{"status": to}
	if to == model.BookingStatusPaid {
		updates["paid_at"] = time.Now()
	}

	if slices.Contains([]model.BookingStatus{
		model.BookingStatusPaid,
		model.BookingStatusExpired,
		model.BookingStatusCanceled,
	}, to) {
		updates["expires_at"] = nil
	}

	res := r.db.WithContext(ctx).Model(&model.EventBooking{}).
		Where("id = ? AND status = ?", id, from).
		Updates(updates)
	return res.RowsAffected == 1, res.Error
}

func (r *EventBookingRepository) SetNullRentalOptionID(ctx context.Context, villaPackageID uint) error {
	updates := map[string]any{"villa_package_id": nil}
	return r.db.WithContext(ctx).
		Model(&model.EventBooking{}).
		Where("villa_package_id = ?", villaPackageID).
		Updates(updates).Error
}

func (r *EventBookingRepository) GetBookingConflicts(
	ctx context.Context,
	villaPackageID uint,
	in time.Time,
	out time.Time,
	excludeID uint,
) ([]model.EventBooking, error) {
	var conflicts []model.EventBooking
	db := r.db.WithContext(ctx).
		Where("villa_package_id = ?", villaPackageID).
		Where("status NOT IN ?", []model.BookingStatus{
			model.BookingStatusCanceled,
			model.BookingStatusExpired,
			model.BookingStatusCompleted,
			model.BookingStatusNotShow,
		}).
		Where("start_at <= ? AND end_at >= ?", out, in)

	if excludeID > 0 {
		db = db.Where("id != ?", excludeID)
	}

	if err := db.Find(&conflicts).Error; err != nil {
		return nil, err
	}

	return conflicts, nil
}
