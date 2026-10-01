package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type SpecificationGalleryRepository struct {
	db *gorm.DB
}

func NewSpecificationGalleryRepository(db *gorm.DB) *SpecificationGalleryRepository {
	return &SpecificationGalleryRepository{db: db}
}

func (r *SpecificationGalleryRepository) WithTx(tx *gorm.DB) *SpecificationGalleryRepository {
	return &SpecificationGalleryRepository{db: tx}
}

func (r *SpecificationGalleryRepository) Get(ctx context.Context) (*model.SpecificationGallery, error) {
	var specificationGallery model.SpecificationGallery

	if err := r.db.WithContext(ctx).Model(&model.SpecificationGallery{}).First(&specificationGallery).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &specificationGallery, nil
}

func (r *SpecificationGalleryRepository) Exist(ctx context.Context) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.SpecificationGallery{}).
		Count(&count).Error

	return count > 0, err
}

func (r *SpecificationGalleryRepository) Create(ctx context.Context, specificationGallery *model.SpecificationGallery) error {
	return r.db.WithContext(ctx).Create(specificationGallery).Error
}

func (r *SpecificationGalleryRepository) Update(ctx context.Context, specificationGallery *model.SpecificationGallery) error {
	return r.db.WithContext(ctx).Save(specificationGallery).Error
}

func (r *SpecificationGalleryRepository) Delete(ctx context.Context, specificationGallery *model.SpecificationGallery) error {
	return r.db.WithContext(ctx).Delete(specificationGallery).Error
}
