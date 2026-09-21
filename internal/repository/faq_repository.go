package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type FaqRepository struct {
	db *gorm.DB
}

func NewFaqRepository(db *gorm.DB) *FaqRepository {
	return &FaqRepository{db: db}
}

func (r *FaqRepository) WithTx(tx *gorm.DB) *FaqRepository {
	return &FaqRepository{db: tx}
}

func (r *FaqRepository) GetAll(ctx context.Context) ([]model.Faq, error) {
	var faqs []model.Faq

	if err := r.db.WithContext(ctx).
		Order("sort_order ASC").
		Find(&faqs).Error; err != nil {
		return nil, err
	}

	return faqs, nil
}

func (r *FaqRepository) GetByID(ctx context.Context, id uint) (*model.Faq, error) {
	var faq model.Faq

	if err := r.db.WithContext(ctx).First(&faq, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &faq, nil
}

func (r *FaqRepository) Create(ctx context.Context, faq *model.Faq) error {
	return r.db.WithContext(ctx).Create(faq).Error
}

func (r *FaqRepository) Update(ctx context.Context, faq *model.Faq) error {
	return r.db.WithContext(ctx).Save(faq).Error
}

func (r *FaqRepository) Delete(ctx context.Context, faq *model.Faq) error {
	return r.db.WithContext(ctx).Delete(faq).Error
}

func (r *FaqRepository) GetMaxSortOrder(ctx context.Context) (int, error) {
	var max int

	err := r.db.WithContext(ctx).
		Model(&model.Faq{}).
		Select("COALESCE(MAX(sort_order), 0)").
		Scan(&max).Error

	return max, err
}

func (r *FaqRepository) UpdateSortOrderRange(ctx context.Context, start, end, delta int) error {
	return r.db.WithContext(ctx).
		Model(&model.Faq{}).
		Where("sort_order BETWEEN ? AND ?", start, end).
		Update("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}

func (r *FaqRepository) ShiftSortOrder(ctx context.Context, sortOrder int) error {
	return r.db.WithContext(ctx).
		Model(&model.Faq{}).
		Where("sort_order >= ?", sortOrder).
		Update("sort_order", gorm.Expr("sort_order + 1")).Error
}
