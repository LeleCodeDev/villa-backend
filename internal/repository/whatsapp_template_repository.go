package repository

import (
	"context"

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
	var templates []model.WhatsappTemplate
	var total int64

	db := r.db.WithContext(ctx).Model(&model.WhatsappTemplate{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&templates).Error; err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}
