package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type WhatsappTemplateRepository struct {
	db *gorm.DB
}

func NewWhatsappTemplateRepository(db *gorm.DB) *WhatsappTemplateRepository {
	return &WhatsappTemplateRepository{db: db}
}

func (r *WhatsappTemplateRepository) WithTx(tx *gorm.DB) *WhatsappTemplateRepository {
	return &WhatsappTemplateRepository{db: tx}
}

func (r *WhatsappTemplateRepository) GetAll(ctx context.Context, query dto.WhatsappTemplateQuery) ([]model.WhatsappTemplate, int64, error) {
	var whatsappTemplates []model.WhatsappTemplate
	var total int64

	db := r.db.WithContext(ctx).Model(&model.WhatsappTemplate{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&whatsappTemplates).Error; err != nil {
		return nil, 0, err
	}

	return whatsappTemplates, total, nil
}

func (r *WhatsappTemplateRepository) GetByID(ctx context.Context, id uint) (*model.WhatsappTemplate, error) {
	var whatsappTemplate model.WhatsappTemplate

	if err := r.db.WithContext(ctx).First(&whatsappTemplate, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &whatsappTemplate, nil
}

func (r *WhatsappTemplateRepository) Create(ctx context.Context, whatsappTemplate *model.WhatsappTemplate) error {
	return r.db.WithContext(ctx).Create(whatsappTemplate).Error
}

func (r *WhatsappTemplateRepository) Update(ctx context.Context, whatsappTemplate *model.WhatsappTemplate) error {
	return r.db.WithContext(ctx).Save(whatsappTemplate).Error
}

func (r *WhatsappTemplateRepository) Delete(ctx context.Context, whatsappTemplate *model.WhatsappTemplate) error {
	return r.db.WithContext(ctx).Delete(whatsappTemplate).Error
}
