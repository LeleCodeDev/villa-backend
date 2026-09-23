package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) WithTx(tx *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: tx}
}

func (r *ApplicationRepository) Get(ctx context.Context) (*model.Application, error) {
	var application model.Application

	if err := r.db.WithContext(ctx).Model(&model.Application{}).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &application, nil
}

func (r *ApplicationRepository) Exist(ctx context.Context) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Application{}).
		Count(&count).Error

	return count > 0, err
}

func (r *ApplicationRepository) Create(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Create(application).Error
}

func (r *ApplicationRepository) Update(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Save(application).Error
}

func (r *ApplicationRepository) Delete(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Delete(application).Error
}
