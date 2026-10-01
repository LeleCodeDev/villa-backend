package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type AddonListRepository struct {
	db *gorm.DB
}

func NewAddonListRepository(db *gorm.DB) *AddonListRepository {
	return &AddonListRepository{db: db}
}

func (r *AddonListRepository) WithTx(tx *gorm.DB) *AddonListRepository {
	return &AddonListRepository{db: tx}
}

func (r *AddonListRepository) GetAll(ctx context.Context, query dto.AddonListQuery) ([]model.AddonList, int64, error) {
	var addonLists []model.AddonList
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AddonList{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&addonLists).Error; err != nil {
		return nil, 0, err
	}

	return addonLists, total, nil
}

func (r *AddonListRepository) GetAllByAddonIDWithQuery(ctx context.Context, addonID uint, query dto.AddonListQuery) ([]model.AddonList, int64, error) {
	var addonLists []model.AddonList
	var total int64

	db := r.db.WithContext(ctx).Model(&model.AddonList{})

	db = db.Where("addon_id = ?", addonID)

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&addonLists).Error; err != nil {
		return nil, 0, err
	}

	return addonLists, total, nil
}

func (r *AddonListRepository) GetByID(ctx context.Context, id uint) (*model.AddonList, error) {
	var addonList model.AddonList

	if err := r.db.WithContext(ctx).First(&addonList, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &addonList, nil
}

func (r *AddonListRepository) GetAllByAddonIDIn(ctx context.Context, addonIds []uint) ([]model.AddonList, error) {
	var addonLists []model.AddonList

	if err := r.db.WithContext(ctx).
		Where("addon_id IN ?", addonIds).
		Order("sort_order ASC").
		Find(&addonLists).Error; err != nil {
		return nil, err
	}

	return addonLists, nil
}

func (r *AddonListRepository) GetAllByAddonID(ctx context.Context, addonID uint) ([]model.AddonList, error) {
	var addonLists []model.AddonList

	if err := r.db.WithContext(ctx).
		Where("addon_id = ?", addonID).
		Order("sort_order ASC").
		Find(&addonLists).Error; err != nil {
		return nil, err
	}

	return addonLists, nil
}

func (r *AddonListRepository) Create(ctx context.Context, addonList *model.AddonList) error {
	return r.db.WithContext(ctx).Create(addonList).Error
}

func (r *AddonListRepository) Update(ctx context.Context, addonList *model.AddonList) error {
	return r.db.WithContext(ctx).Save(addonList).Error
}

func (r *AddonListRepository) Delete(ctx context.Context, addonList *model.AddonList) error {
	return r.db.WithContext(ctx).Delete(addonList).Error
}

func (r *AddonListRepository) DeleteByAddonID(ctx context.Context, addonID uint) error {
	return r.db.WithContext(ctx).
		Where("addon_id = ?", addonID).
		Delete(&model.AddonList{}).Error
}

func (r *AddonListRepository) GetMaxSortOrder(ctx context.Context, addonID uint) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.AddonList{}).
		Where("addon_id = ?", addonID).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *AddonListRepository) UpdateSortOrderRange(ctx context.Context, addonID uint, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.AddonList{}).
		Where("addon_id = ?", addonID).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *AddonListRepository) ShiftSortOrder(ctx context.Context, addonID uint, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.AddonList{}).
		Where("addon_id = ?", addonID).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
