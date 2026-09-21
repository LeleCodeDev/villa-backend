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

func (s *GalleryService) Update(ctx context.Context, req dto.GalleryRequest, id uint) (dto.GalleryResponse, error) {
	var updatedGallery *model.Gallery

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		gallery, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if gallery == nil {
			return errors.NotFound(fmt.Sprintf("Gallery not found with ID: %d", id))
		}

		oldOrder := gallery.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if newOrder != &oldOrder {
			if *newOrder < oldOrder {
				if err := txRepo.UpdateSortOrderRange(ctx, *newOrder, oldOrder-1, +1); err != nil {
					return err
				}
			} else {
				if err := txRepo.UpdateSortOrderRange(ctx, oldOrder+1, *newOrder, -1); err != nil {
					return err
				}
			}
		}

		imagepath := gallery.Image
		if req.Image != nil {
			if gallery.Image != nil {
				image.DeleteImage(*gallery.Image)
			}

			path, err := image.SaveImage("uploads/galleries", req.Image, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath = &path
		}

		mapper.UpdateGalleryModel(gallery, req, imagepath, *newOrder)
		if err := txRepo.Update(ctx, gallery); err != nil {
			image.DeleteImage(*gallery.Image)
			return err
		}

		updatedGallery = gallery

		return nil
	}); err != nil {
		return dto.GalleryResponse{}, err
	}

	return mapper.ToGalleryResponse(updatedGallery), nil
}

func (s *GalleryService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		gallery, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if gallery == nil {
			return errors.NotFound(fmt.Sprintf("Gallery not found with ID: %d", id))
		}

		if gallery.Image != nil {
			image.DeleteImage(*gallery.Image)
		}

		return txRepo.Delete(ctx, gallery)
	})
}
