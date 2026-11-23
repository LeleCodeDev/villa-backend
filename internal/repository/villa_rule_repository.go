package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type VillaRuleRepository struct {
	db *gorm.DB
}

func NewVillaRuleRepository(db *gorm.DB) *VillaRuleRepository {
	return &VillaRuleRepository{db: db}
}

func (r *VillaRuleRepository) WithTx(tx *gorm.DB) *VillaRuleRepository {
	return &VillaRuleRepository{db: tx}
}

func (r *VillaRuleRepository) GetAll(ctx context.Context, query dto.VillaRuleQuery) ([]model.VillaRule, int64, error) {
	var villaRules []model.VillaRule
	var total int64

	db := r.db.WithContext(ctx).Model(&model.VillaRule{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("sort_order " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&villaRules).Error; err != nil {
		return nil, 0, err
	}

	return villaRules, total, nil
}

func (r *VillaRuleRepository) GetByID(ctx context.Context, id uint) (*model.VillaRule, error) {
	var villaRule model.VillaRule

	if err := r.db.WithContext(ctx).First(&villaRule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &villaRule, nil
}

func (r *VillaRuleRepository) Create(ctx context.Context, villaRule *model.VillaRule) error {
	return r.db.WithContext(ctx).Create(villaRule).Error
}

func (r *VillaRuleRepository) Update(ctx context.Context, villaRule *model.VillaRule) error {
	return r.db.WithContext(ctx).Save(villaRule).Error
}

func (r *VillaRuleRepository) Delete(ctx context.Context, villaRule *model.VillaRule) error {
	return r.db.WithContext(ctx).Delete(villaRule).Error
}

func (r *VillaRuleRepository) ExistByLogoID(ctx context.Context, logoID uint) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.VillaRule{}).
		Where("logo_id = ?", logoID).
		Count(&count).Error

	return count > 0, err
}

func (r *VillaRuleRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.VillaRule{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *VillaRuleRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaRule{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *VillaRuleRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.VillaRule{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
