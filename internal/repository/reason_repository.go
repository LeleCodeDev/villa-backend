package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type ReasonRepository struct {
	db *gorm.DB
}

func NewReasonRepository(db *gorm.DB) *ReasonRepository {
	return &ReasonRepository{
		db: db,
	}
}

func (r *ReasonRepository) WithTx(tx *gorm.DB) *ReasonRepository {
	return &ReasonRepository{db: tx}
}

func (r *ReasonRepository) GetAll(ctx context.Context, query dto.ReasonQuery) ([]model.Reason, int64, error) {
	var reasons []model.Reason
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Reason{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&reasons).Error; err != nil {
		return nil, 0, err
	}

	return reasons, total, nil
}

func (r *ReasonRepository) GetByID(ctx context.Context, id uint) (*model.Reason, error) {
	var reason model.Reason

	if err := r.db.WithContext(ctx).First(&reason, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &reason, nil
}

func (r *ReasonRepository) Create(ctx context.Context, reason *model.Reason) error {
	return r.db.WithContext(ctx).Create(reason).Error
}

func (r *ReasonRepository) Update(ctx context.Context, reason *model.Reason) error {
	return r.db.WithContext(ctx).Save(reason).Error
}

func (r *ReasonRepository) Delete(ctx context.Context, reason *model.Reason) error {
	return r.db.WithContext(ctx).Delete(reason).Error
}
