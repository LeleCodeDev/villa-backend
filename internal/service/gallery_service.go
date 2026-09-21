package service

import (
	"context"
	"fmt"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"github.com/lelecodedev/villa-backend/pkg/image"
	"gorm.io/gorm"
)

type GalleryService struct {
	txManager *repository.TxManager
	repo      *repository.GalleryRepository
}

func NewGalleryService(
	txManager *repository.TxManager,
	repo *repository.GalleryRepository,
) *GalleryService {
	return &GalleryService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *GalleryService) GetAll(ctx context.Context) ([]dto.GalleryResponse, error) {
	galleries, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GalleryResponse, 0, len(galleries))
	for _, gallery := range galleries {
		responses = append(responses, mapper.ToGalleryResponse(&gallery))
	}

	return responses, nil
}

func (s *GalleryService) GetByID(ctx context.Context, id uint) (dto.GalleryResponse, error) {
	gallery, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.GalleryResponse{}, err
	}
	if gallery == nil {
		return dto.GalleryResponse{}, errors.NotFound(fmt.Sprintf("Gallery not found with ID: %d", id))
	}

	return mapper.ToGalleryResponse(gallery), nil
}

func (s *GalleryService) Create(ctx context.Context, req dto.GalleryRequest) (dto.GalleryResponse, error) {
	var createdGallery *model.Gallery

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		var sortOrder int
		if req.SortOrder == nil {
			maxOrder, err := txRepo.GetMaxSortOrder(ctx)
			if err != nil {
				return err
			}

			sortOrder = maxOrder + 1
		} else {
			sortOrder = *req.SortOrder
			if err := txRepo.ShiftSortOrder(ctx, sortOrder); err != nil {
				return err
			}
		}

		var imagePath *string
		if req.Image != nil {
			path, err := image.SaveImage("uploads/galleries", req.Image, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath = &path
		}

		gallery := mapper.ToGalleryModel(req, imagePath, sortOrder)
		if err := txRepo.Create(ctx, gallery); err != nil {
			image.DeleteImage(*imagePath)
			return err
		}

		createdGallery = gallery

		return nil
	}); err != nil {
		return dto.GalleryResponse{}, err
	}

	return mapper.ToGalleryResponse(createdGallery), nil
}
