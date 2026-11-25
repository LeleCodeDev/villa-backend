package service

import (
	"context"
	"log"
	"mime/multipart"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/media"
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
	var saved []string

	err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist("Specification gallery data already exist!")
		}

		var paths [3]*string
		files := [3]*multipart.FileHeader{req.Image1, req.Image2, req.Image3}

		for i, file := range files {
			if file == nil {
				continue
			}

			path, err := image.Save("uploads/specification_galleries", file, image.DefaultOptions())
			if err != nil {
				return err
			}

			saved = append(saved, path)
			paths[i] = &path
		}

		specificationGallery := mapper.ToSpecificationGalleryModel(paths[0], paths[1], paths[2])
		if err := txRepo.Create(ctx, specificationGallery); err != nil {
			return err
		}

		createdSpecificationGallery = specificationGallery
		return nil
	})
	if err != nil {
		for _, imagePath := range saved {
			if delErr := media.Delete(imagePath); delErr != nil {
				log.Printf("failed to cleanup image%s: %v", imagePath, delErr)
			}
		}
		return dto.SpecificationGalleryResponse{}, err
	}

	return mapper.ToSpecificationGalleryResponse(createdSpecificationGallery), nil
}

func (s *SpecificationGalleryService) Update(ctx context.Context, req dto.SpecificationGalleryRequest) (dto.SpecificationGalleryResponse, error) {
	var updatedSpecificationGallery *model.SpecificationGallery
	var saved []string
	var replaced []string

	err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		specificationGallery, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if specificationGallery == nil {
			return appError.NotFound("Specification gallery data not found!")
		}

		current := [3]*string{specificationGallery.Image1, specificationGallery.Image2, specificationGallery.Image3}
		files := [3]*multipart.FileHeader{req.Image1, req.Image2, req.Image3}
		paths := current

		for i, file := range files {
			if file == nil {
				continue
			}

			path, err := image.Save("uploads/specification_galleries", file, image.DefaultOptions())
			if err != nil {
				return err
			}

			saved = append(saved, path)
			if current[i] != nil {
				replaced = append(replaced, *current[i])
			}
			paths[i] = &path
		}

		mapper.UpdateSpecificationGalleryModel(specificationGallery, paths[0], paths[1], paths[2])
		if err := txRepo.Update(ctx, specificationGallery); err != nil {
			return err
		}

		updatedSpecificationGallery = specificationGallery
		return nil
	})
	if err != nil {
		for _, imagePath := range saved {
			if delErr := media.Delete(imagePath); delErr != nil {
				log.Printf("failed to cleanup new image %s: %v", imagePath, delErr)
			}
		}
		return dto.SpecificationGalleryResponse{}, err
	}

	for _, imagePath := range replaced {
		if delErr := media.Delete(imagePath); delErr != nil {
			log.Printf("failed to delete old image %s: %v", imagePath, delErr)
		}
	}

	return mapper.ToSpecificationGalleryResponse(updatedSpecificationGallery), nil
}

func (s *SpecificationGalleryService) Delete(ctx context.Context) error {
	var imagePaths []string

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		specificationGallery, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if specificationGallery == nil {
			return appError.NotFound("Specification gallery data not found!")
		}

		for _, path := range []*string{specificationGallery.Image1, specificationGallery.Image2, specificationGallery.Image3} {
			if path != nil {
				imagePaths = append(imagePaths, *path)
			}
		}

		return txRepo.Delete(ctx, specificationGallery)
	}); err != nil {
		return err
	}

	for _, imagePath := range imagePaths {
		if delErr := media.Delete(imagePath); delErr != nil {
			log.Printf("failed to cleanup image %s: %v", imagePath, delErr)
		}
	}

	return nil
}
