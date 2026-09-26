package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type VillaPackageRepository struct {
	db *gorm.DB
}

func NewVillaPackageRepository(db *gorm.DB) *VillaPackageRepository {
	return &VillaPackageRepository{db: db}
}

func (r *VillaPackageRepository) WithTx(tx *gorm.DB) *VillaPackageRepository {
	return &VillaPackageRepository{db: tx}
}

func (r *VillaPackageRepository) GetAll(ctx context.Context) ([]model.VillaPackage, error) {
	var villaPackages []model.VillaPackage

	if err := r.db.WithContext(ctx).
		Order("sort_order ASC").
		Find(&villaPackages).Error; err != nil {
		return nil, err
	}

	return villaPackages, nil
}

func (r *VillaPackageRepository) GetByID(ctx context.Context, id uint) (*model.VillaPackage, error) {
	var villaPackage model.VillaPackage

	if err := r.db.WithContext(ctx).First(&villaPackage, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &villaPackage, nil
}

func (r *VillaPackageRepository) ExistByID(ctx context.Context, id uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.VillaPackage{}).
		Where("id = ?", id).
		Count(&count).Error

	return count > 0, err
}

func (r *VillaPackageRepository) Create(ctx context.Context, villaPackage *model.VillaPackage) error {
	return r.db.WithContext(ctx).Create(villaPackage).Error
}

func (r *VillaPackageRepository) Update(ctx context.Context, villaPackage *model.VillaPackage) error {
	return r.db.WithContext(ctx).Save(villaPackage).Error
}

func (r *VillaPackageRepository) Delete(ctx context.Context, villaPackage *model.VillaPackage) error {
	return r.db.WithContext(ctx).Delete(villaPackage).Error
}

func (r *VillaPackageRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.VillaPackage{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *VillaPackageRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaPackage{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *VillaPackageRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaPackage{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
