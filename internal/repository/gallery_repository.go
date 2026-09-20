package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type GalleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
	return &GalleryRepository{db: db}
}

func (r *GalleryRepository) WithTx(tx *gorm.DB) *GalleryRepository {
	return &GalleryRepository{db: tx}
}

func (r *GalleryRepository) GetAll(ctx context.Context) ([]model.Gallery, error) {
	var galleries []model.Gallery

	if err := r.db.WithContext(ctx).Find(&galleries).Error; err != nil {
		return nil, err
	}

	return galleries, nil
}

func (r *GalleryRepository) GetByID(ctx context.Context, id uint) (*model.Gallery, error) {
	var gallery model.Gallery

	if err := r.db.WithContext(ctx).First(&gallery, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &gallery, nil
}

func (r *GalleryRepository) Create(ctx context.Context, gallery *model.Gallery) error {
	return r.db.WithContext(ctx).Create(gallery).Error
}

func (r *GalleryRepository) Update(ctx context.Context, gallery *model.Gallery) error {
	return r.db.WithContext(ctx).Save(gallery).Error
}

func (r *GalleryRepository) Delete(ctx context.Context, gallery *model.Gallery) error {
	return r.db.WithContext(ctx).Delete(gallery).Error
}

func (r *GalleryRepository) GetMaxOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Select("COALESCE(MAX(order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *GalleryRepository) UpdateOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Where("order BETWEEN ? AND ?", start, end).
		Update("order", gorm.Expr("order + ", delta)).Error
}

func (r *GalleryRepository) ShiftOrder(ctx context.Context, order int) error {
	return r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Where("order >= ?", order).
		Update("order", gorm.Expr("order + 1")).Error
}
