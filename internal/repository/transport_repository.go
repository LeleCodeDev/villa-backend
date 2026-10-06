package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type TransportRepository struct {
	db *gorm.DB
}

func NewTransportRepository(db *gorm.DB) *TransportRepository {
	return &TransportRepository{db: db}
}

func (r *TransportRepository) WithTx(tx *gorm.DB) *TransportRepository {
	return &TransportRepository{db: tx}
}

func (r *TransportRepository) GetAll(ctx context.Context, query dto.TransportQuery) ([]model.Transport, int64, error) {
	var transports []model.Transport
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Transport{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&transports).Error; err != nil {
		return nil, 0, err
	}

	return transports, total, nil
}

func (r *TransportRepository) GetAllByAttractionIDWithQuery(ctx context.Context, attractionID uint, query dto.TransportQuery) ([]model.Transport, int64, error) {
	var transports []model.Transport
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Transport{})

	db = db.Where("attraction_id = ?", attractionID)

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&transports).Error; err != nil {
		return nil, 0, err
	}

	return transports, total, nil
}

func (r *TransportRepository) GetByID(ctx context.Context, id uint) (*model.Transport, error) {
	var transport model.Transport

	if err := r.db.WithContext(ctx).First(&transport, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &transport, nil
}

func (r *TransportRepository) GetAllByAttractionIDIn(ctx context.Context, attractionIds []uint) ([]model.Transport, error) {
	var transports []model.Transport

	if err := r.db.WithContext(ctx).
		Where("attraction_id IN ?", attractionIds).
		Order("sort_order ASC").
		Find(&transports).Error; err != nil {
		return nil, err
	}

	return transports, nil
}

func (r *TransportRepository) GetAllByAttractionID(ctx context.Context, attractionID uint) ([]model.Transport, error) {
	var transports []model.Transport

	if err := r.db.WithContext(ctx).
		Where("attraction_id = ?", attractionID).
		Order("sort_order ASC").
		Find(&transports).Error; err != nil {
		return nil, err
	}

	return transports, nil
}

func (r *TransportRepository) Create(ctx context.Context, transport *model.Transport) error {
	return r.db.WithContext(ctx).Create(transport).Error
}

func (r *TransportRepository) Update(ctx context.Context, transport *model.Transport) error {
	return r.db.WithContext(ctx).Save(transport).Error
}

func (r *TransportRepository) Delete(ctx context.Context, transport *model.Transport) error {
	return r.db.WithContext(ctx).Delete(transport).Error
}

func (r *TransportRepository) DeleteByAttractionID(ctx context.Context, attractionID uint) error {
	return r.db.WithContext(ctx).
		Where("attraction_id = ?", attractionID).
		Delete(&model.Transport{}).Error
}

func (r *TransportRepository) GetMaxSortOrder(ctx context.Context, attractionID uint) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Transport{}).
		Where("attraction_id = ?", attractionID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *TransportRepository) UpdateSortOrderRange(ctx context.Context, attractionID uint, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Transport{}).
		Where("attraction_id = ?", attractionID).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *TransportRepository) ShiftSortOrder(ctx context.Context, attractionID uint, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Transport{}).
		Where("attraction_id = ?", attractionID).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
