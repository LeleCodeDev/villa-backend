package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"github.com/lelecodedev/villa-backend/pkg/image"
	"gorm.io/gorm"
)

type ApplicationService struct {
	txManager *repository.TxManager
	repo      *repository.ApplicationRepository
}

func NewApplicationService(
	txManager *repository.TxManager,
	repo *repository.ApplicationRepository,
) *ApplicationService {
	return &ApplicationService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *ApplicationService) Get(ctx context.Context) (dto.ApplicationResponse, error) {
	application, err := s.repo.Get(ctx)
	if err != nil {
		return dto.ApplicationResponse{}, err
	}
	if application == nil {
		return dto.ApplicationResponse{}, errors.NotFound("Application data not found!")
	}

	return mapper.ToApplicationResponse(application), nil
}

func (s *ApplicationService) Create(ctx context.Context, req dto.ApplicationRequest) (dto.ApplicationResponse, error) {
	var createdApplication *model.Application

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			return errors.AlreadyExist("Application data already exist!")
		}

		var imagePath *string
		if req.Logo != nil {
			path, err := image.SaveImage("uploads/application", req.Logo, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath = &path
		}

		application := mapper.ToApplicationModel(req, imagePath)
		if err := txRepo.Create(ctx, application); err != nil {
			image.DeleteImage(*imagePath)
			return err
		}

		createdApplication = application

		return nil
	}); err != nil {
		return dto.ApplicationResponse{}, err
	}

	return mapper.ToApplicationResponse(createdApplication), nil
}

func (s *ApplicationService) Update(ctx context.Context, req dto.ApplicationRequest) (dto.ApplicationResponse, error) {
	var updatedApplication *model.Application

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		application, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if application == nil {
			return errors.NotFound("Application data not found!")
		}

		imagepath := application.Logo
		if req.Logo != nil {
			if application.Logo != nil {
				image.DeleteImage(*application.Logo)
			}

			path, err := image.SaveImage("uploads/application", req.Logo, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath = &path
		}

		mapper.UpdateApplicationModel(application, req, imagepath)
		if err := txRepo.Update(ctx, application); err != nil {
			image.DeleteImage(*application.Logo)
			return err
		}

		updatedApplication = application

		return nil
	}); err != nil {
		return dto.ApplicationResponse{}, err
	}

	return mapper.ToApplicationResponse(updatedApplication), nil
}

func (s *ApplicationService) Delete(ctx context.Context) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		application, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if application == nil {
			return errors.NotFound("Application data not found!")
		}

		return txRepo.Delete(ctx, application)
	})
}
