package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type WhatsappTemplateService struct {
	txManager *repository.TxManager
	repo      *repository.WhatsappTemplateRepository
}

func NewWhatsappTemplateService(
	txManager *repository.TxManager,
	repo *repository.WhatsappTemplateRepository,
) *WhatsappTemplateService {
	return &WhatsappTemplateService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *WhatsappTemplateService) GetAll(ctx context.Context, query dto.WhatsappTemplateQuery) ([]dto.WhatsappTemplateResponse, int64, error) {
	whatsappTemplates, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	whatsappTemplateResponses := make([]dto.WhatsappTemplateResponse, 0, len(whatsappTemplates))
	for _, whatsappTemplate := range whatsappTemplates {
		whatsappTemplateResponses = append(whatsappTemplateResponses, mapper.ToWhatsappTemplateResponse(&whatsappTemplate))
	}

	return whatsappTemplateResponses, total, nil
}

func (s *WhatsappTemplateService) GetByID(ctx context.Context, id uint) (dto.WhatsappTemplateResponse, error) {
	whatsappTemplate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.WhatsappTemplateResponse{}, err
	}
	if whatsappTemplate == nil {
		return dto.WhatsappTemplateResponse{}, appError.BadRequest(fmt.Sprintf("WhatsappTemplate not found with ID: %d", id))
	}

	return mapper.ToWhatsappTemplateResponse(whatsappTemplate), nil
}

func (s *WhatsappTemplateService) Create(ctx context.Context, req dto.WhatsappTemplateRequest) (dto.WhatsappTemplateResponse, error) {
	var createdWhatsappTemplate *model.WhatsappTemplate

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		whatsappTemplate := mapper.ToWhatsappTemplateModel(req)
		if err := txRepo.Create(ctx, whatsappTemplate); err != nil {
			return err
		}

		createdWhatsappTemplate = whatsappTemplate

		return nil
	}); err != nil {
		return dto.WhatsappTemplateResponse{}, err
	}

	return mapper.ToWhatsappTemplateResponse(createdWhatsappTemplate), nil
}

func (s *WhatsappTemplateService) Update(ctx context.Context, req dto.WhatsappTemplateRequest, id uint) (dto.WhatsappTemplateResponse, error) {
	var updatedWhatsappTemplate *model.WhatsappTemplate

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		whatsappTemplate, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if whatsappTemplate == nil {
			return appError.NotFound(fmt.Sprintf("WhatsappTemplate not found with ID: %d", id))
		}

		mapper.UpdateWhatsappTemplateModel(whatsappTemplate, req)
		if err := txRepo.Update(ctx, whatsappTemplate); err != nil {
			return err
		}

		updatedWhatsappTemplate = whatsappTemplate

		return nil
	}); err != nil {
		return dto.WhatsappTemplateResponse{}, err
	}

	return mapper.ToWhatsappTemplateResponse(updatedWhatsappTemplate), nil
}

func (s *WhatsappTemplateService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		whatsappTemplate, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if whatsappTemplate == nil {
			return appError.NotFound(fmt.Sprintf("WhatsappTemplate not found with ID: %d", id))
		}

		return txRepo.Delete(ctx, whatsappTemplate)
	})
}
