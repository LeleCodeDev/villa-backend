package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type VillaPackageListRepository struct {
	db *gorm.DB
}

func NewVillaPackageListRepository(db *gorm.DB) *VillaPackageListRepository {
	return &VillaPackageListRepository{db: db}
}

func (r *VillaPackageListRepository) WithTx(tx *gorm.DB) *VillaPackageListRepository {
	return &VillaPackageListRepository{db: tx}
}

func (r *VillaPackageListRepository) GetAll(ctx context.Context) ([]model.VillaPackageList, error) {
	var villaPackageLists []model.VillaPackageList

	if err := r.db.WithContext(ctx).
		Order("sort_order ASC").
		Find(&villaPackageLists).Error; err != nil {
		return nil, err
	}

	return villaPackageLists, nil
}

func (r *VillaPackageListRepository) GetByID(ctx context.Context, id uint) (*model.VillaPackageList, error) {
	var villaPackageList model.VillaPackageList

	if err := r.db.WithContext(ctx).First(&villaPackageList, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &villaPackageList, nil
}

func (r *VillaPackageListRepository) GetAllByVillaPackageIDIn(ctx context.Context, villaPackageIds []uint) ([]model.VillaPackageList, error) {
	var villaPackageLists []model.VillaPackageList

	if err := r.db.WithContext(ctx).
		Where("villa_package_id IN ?", villaPackageIds).
		Order("sort_order ASC").
		Find(&villaPackageLists).Error; err != nil {
		return nil, err
	}

	return villaPackageLists, nil
}

func (r *VillaPackageListRepository) GetAllByVillaPackageID(ctx context.Context, villaPackageId uint) ([]model.VillaPackageList, error) {
	var villaPackageLists []model.VillaPackageList

	if err := r.db.WithContext(ctx).
		Where("villa_package_id = ?", villaPackageId).
		Order("sort_order ASC").
		Find(&villaPackageLists).Error; err != nil {
		return nil, err
	}

	return villaPackageLists, nil
}

func (r *VillaPackageListRepository) Create(ctx context.Context, villaPackageList *model.VillaPackageList) error {
	return r.db.WithContext(ctx).Create(villaPackageList).Error
}

func (r *VillaPackageListRepository) Update(ctx context.Context, villaPackageList *model.VillaPackageList) error {
	return r.db.WithContext(ctx).Save(villaPackageList).Error
}

func (r *VillaPackageListRepository) Delete(ctx context.Context, villaPackageList *model.VillaPackageList) error {
	return r.db.WithContext(ctx).Delete(villaPackageList).Error
}

func (r *VillaPackageListRepository) DeleteByVillaPackageID(ctx context.Context, villaPackageId uint) error {
	return r.db.WithContext(ctx).
		Where("villa_package_id = ?", villaPackageId).
		Delete(&model.VillaPackageList{}).Error
}

func (r *VillaPackageListRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.VillaPackageList{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *VillaPackageListRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaPackageList{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *VillaPackageListRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaPackageList{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
