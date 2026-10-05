package service

import (
	"context"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"

	"github.com/lelecodedev/villa-backend/internal/image"
	"gorm.io/gorm"
)

type SpecificationGalleryService struct {
	txManager *repository.TxManager
	repo      *repository.SpecificationGalleryRepository
}

func NewSpecificationGalleryService(
	txManager *repository.TxManager,
	repo *repository.SpecificationGalleryRepository,
) *SpecificationGalleryService {
	return &SpecificationGalleryService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *SpecificationGalleryService) Get(ctx context.Context) (dto.SpecificationGalleryResponse, error) {
	specificationGallery, err := s.repo.Get(ctx)
	if err != nil {
		return dto.SpecificationGalleryResponse{}, err
	}
	if specificationGallery == nil {
		return dto.SpecificationGalleryResponse{}, appError.NotFound("Specification gallery data not found!")
	}

	return mapper.ToSpecificationGalleryResponse(specificationGallery), nil
}

func (s *SpecificationGalleryService) Create(ctx context.Context, req dto.SpecificationGalleryRequest) (dto.SpecificationGalleryResponse, error) {
	var createdSpecificationGallery *model.SpecificationGallery

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist("Specification gallery data already exist!")
		}

		var imagePath1 *string
		if req.Image1 != nil {
			path, err := image.SaveImage("uploads/specification_galleries", req.Image1, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath1 = &path
		}

		var imagePath2 *string
		if req.Image2 != nil {
			path, err := image.SaveImage("uploads/specification_galleries", req.Image2, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath2 = &path
		}

		var imagePath3 *string
		if req.Image3 != nil {
			path, err := image.SaveImage("uploads/specification_galleries", req.Image3, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath3 = &path
		}

		specificationGallery := mapper.ToSpecificationGalleryModel(imagePath1, imagePath2, imagePath3)
		if err := txRepo.Create(ctx, specificationGallery); err != nil {
			image.DeleteImage(*imagePath1)
			image.DeleteImage(*imagePath2)
			image.DeleteImage(*imagePath3)
			return err
		}

		createdSpecificationGallery = specificationGallery

		return nil
	}); err != nil {
		return dto.SpecificationGalleryResponse{}, err
	}

	return mapper.ToSpecificationGalleryResponse(createdSpecificationGallery), nil
}

func (s *SpecificationGalleryService) Update(ctx context.Context, req dto.SpecificationGalleryRequest) (dto.SpecificationGalleryResponse, error) {
	var updatedSpecificationGallery *model.SpecificationGallery

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		specificationGallery, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if specificationGallery == nil {
			return appError.NotFound("Specification gallery data not found!")
		}

		imagepath1 := specificationGallery.Image1
		if req.Image1 != nil {
			if specificationGallery.Image1 != nil {
				image.DeleteImage(*specificationGallery.Image1)
			}

			path, err := image.SaveImage("uploads/galleries", req.Image1, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath1 = &path
		}

		imagepath2 := specificationGallery.Image2
		if req.Image2 != nil {
			if specificationGallery.Image2 != nil {
				image.DeleteImage(*specificationGallery.Image2)
			}

			path, err := image.SaveImage("uploads/galleries", req.Image2, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath2 = &path
		}

		imagepath3 := specificationGallery.Image3
		if req.Image3 != nil {
			if specificationGallery.Image3 != nil {
				image.DeleteImage(*specificationGallery.Image3)
			}

			path, err := image.SaveImage("uploads/galleries", req.Image3, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath3 = &path
		}

		mapper.UpdateSpecificationGalleryModel(specificationGallery, imagepath1, imagepath2, imagepath3)
		if err := txRepo.Update(ctx, specificationGallery); err != nil {
			return err
		}

		updatedSpecificationGallery = specificationGallery

		return nil
	}); err != nil {
		return dto.SpecificationGalleryResponse{}, err
	}

	return mapper.ToSpecificationGalleryResponse(updatedSpecificationGallery), nil
}

func (s *SpecificationGalleryService) Delete(ctx context.Context) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		specificationGallery, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if specificationGallery == nil {
			return appError.NotFound("Specification gallery data not found!")
		}

		if specificationGallery.Image1 != nil {
			image.DeleteImage(*specificationGallery.Image1)
		}
		if specificationGallery.Image2 != nil {
			image.DeleteImage(*specificationGallery.Image2)
		}
		if specificationGallery.Image3 != nil {
			image.DeleteImage(*specificationGallery.Image3)
		}

		return txRepo.Delete(ctx, specificationGallery)
	})
}
