package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
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

func (r *GalleryRepository) GetAll(ctx context.Context, query dto.GalleryQuery) ([]model.Gallery, int64, error) {
	var galleries []model.Gallery
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Gallery{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&galleries).Error; err != nil {
		return nil, 0, err
	}

	return galleries, total, nil
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

func (r *GalleryRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *GalleryRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *GalleryRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Gallery{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
