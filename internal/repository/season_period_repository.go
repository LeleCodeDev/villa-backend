package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type SeasonPeriodRepository struct {
	db *gorm.DB
}

func NewSeasonPeriodRepository(db *gorm.DB) *SeasonPeriodRepository {
	return &SeasonPeriodRepository{db: db}
}

func (r *SeasonPeriodRepository) WithTx(tx *gorm.DB) *SeasonPeriodRepository {
	return &SeasonPeriodRepository{db: tx}
}

func (r *SeasonPeriodRepository) GetAll(ctx context.Context, query dto.SeasonPeriodQuery) ([]model.SeasonPeriod, int64, error) {
	var seasonPeriods []model.SeasonPeriod
	var total int64

	db := r.db.WithContext(ctx).Model(&model.SeasonPeriod{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("start_date " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&seasonPeriods).Error; err != nil {
		return nil, 0, err
	}

	return seasonPeriods, total, nil
}

func (r *SeasonPeriodRepository) GetByID(ctx context.Context, id uint) (*model.SeasonPeriod, error) {
	var seasonPeriod model.SeasonPeriod

	if err := r.db.WithContext(ctx).First(&seasonPeriod, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &seasonPeriod, nil
}

func (r *SeasonPeriodRepository) GetSeasonConflict(ctx context.Context, startDate time.Time, endDate time.Time, excludeID uint) ([]model.SeasonPeriod, error) {
	var conflicts []model.SeasonPeriod

	if err := r.db.WithContext(ctx).
		Model(&model.SeasonPeriod{}).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Find(&conflicts).Error; err != nil {
		return nil, err
	}

	return conflicts, nil
}

func (r *SeasonPeriodRepository) Create(ctx context.Context, seasonPeriod *model.SeasonPeriod) error {
	return r.db.WithContext(ctx).Create(seasonPeriod).Error
}

func (r *SeasonPeriodRepository) Update(ctx context.Context, seasonPeriod *model.SeasonPeriod) error {
	return r.db.WithContext(ctx).Save(seasonPeriod).Error
}

func (r *SeasonPeriodRepository) Delete(ctx context.Context, seasonPeriod *model.SeasonPeriod) error {
	return r.db.WithContext(ctx).Delete(seasonPeriod).Error
}
