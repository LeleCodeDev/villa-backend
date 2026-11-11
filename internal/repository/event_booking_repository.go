package repository

import (
	"context"
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

func (r EventBookingRepository) GetAll(ctx context.Context, query dto.EventBookingQuery, userID *uint) ([]model.EventBooking, int64, error) {
	var eventBookings []model.EventBooking
	var total int64

	db := r.db.WithContext(ctx).
		Preload("User").
		Model(&model.EventBooking{})

	if userID != nil {
		db = db.Where("user_id = ?", userID)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if query.StartAt != "" {
		start, err := time.ParseInLocation("2006-01-02 15:04", query.StartAt, config.Env.Location)
		if err != nil {
			return nil, 0, err
		}

		db = db.Where("start_at >= ?", start)
	}

	if query.EndAt != "" {
		end, err := time.ParseInLocation("2006-01-02 15:04", query.EndAt, config.Env.Location)
		if err != nil {
			return nil, 0, err
		}

		db = db.Where("end_at <= ?", end)
	}

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at" + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&eventBookings).Error; err != nil {
		return nil, 0, err
	}

	return eventBookings, total, nil
}
