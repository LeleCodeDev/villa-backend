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

	if err := r.db.WithContext(ctx).Find(&faqs).Error; err != nil {
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
