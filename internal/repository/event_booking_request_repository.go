package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type EventBookingRequestRepository struct {
	db *gorm.DB
}

func NewEventBookingRequestRepository(db *gorm.DB) *EventBookingRequestRepository {
	return &EventBookingRequestRepository{db: db}
}

func (r *EventBookingRequestRepository) WithTx(tx *gorm.DB) *EventBookingRequestRepository {
	return &EventBookingRequestRepository{db: tx}
}

func (r *EventBookingRequestRepository) GetAll(ctx context.Context, query dto.EventBookingRequestQuery, userID *uint) ([]model.EventBookingRequest, int64, error) {
	var eventBookingRequests []model.EventBookingRequest
	var total int64

	db := r.db.WithContext(ctx).
		Preload("User").
		Preload("VillaPackage").
		Model(&model.EventBookingRequest{})

	if userID != nil {
		db = db.Where("user_id = ?", userID)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
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

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&eventBookingRequests).Error; err != nil {
		return nil, 0, err
	}

	return eventBookingRequests, total, nil
}

func (r *EventBookingRequestRepository) GetByID(ctx context.Context, id uint) (*model.EventBookingRequest, error) {
	var request model.EventBookingRequest

	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("VillaPackage").
		First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &request, nil
}

func (r *EventBookingRequestRepository) SetNullVillaPackageID(ctx context.Context, villaPackageID uint) error {
	updates := map[string]any{"villa_package_id": nil}
	return r.db.WithContext(ctx).
		Model(&model.EventBookingRequest{}).
		Where("villa_package_id = ?", villaPackageID).
		Updates(updates).Error
}

func (r *EventBookingRequestRepository) Create(ctx context.Context, request *model.EventBookingRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *EventBookingRequestRepository) Update(ctx context.Context, request *model.EventBookingRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

func (r *EventBookingRequestRepository) Delete(ctx context.Context, request *model.EventBookingRequest) error {
	return r.db.WithContext(ctx).Delete(request).Error
}
