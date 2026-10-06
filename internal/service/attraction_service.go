package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type AttractionService struct {
	txManager     *repository.TxManager
	repo          *repository.AttractionRepository
	transportRepo *repository.TransportRepository
}

func NewAttractionService(
	txManager *repository.TxManager,
	repo *repository.AttractionRepository,
	transportRepo *repository.TransportRepository,
) *AttractionService {
	return &AttractionService{
		txManager:     txManager,
		repo:          repo,
		transportRepo: transportRepo,
	}
}

func (s *AttractionService) GetAll(ctx context.Context, query dto.AttractionQuery) ([]dto.AttractionResponse, int64, error) {
	attractions, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	attractionIds := make([]uint, 0, len(attractions))
	for _, attraction := range attractions {
		attractionIds = append(attractionIds, attraction.ID)
	}

	alltransports, err := s.transportRepo.GetAllByAttractionIDIn(ctx, attractionIds)
	if err != nil {
		return nil, 0, err
	}

	transportMap := make(map[uint][]model.Transport, len(attractions))
	for _, transport := range alltransports {
		transportMap[transport.AttractionID] = append(transportMap[transport.AttractionID], transport)
	}

	responses := make([]dto.AttractionResponse, 0, len(attractions))
	for _, attraction := range attractions {
		responses = append(responses, mapper.ToAttractionResponse(&attraction, transportMap[attraction.ID]))
	}

	return responses, total, nil
}

func (s *AttractionService) GetByID(ctx context.Context, id uint) (dto.AttractionResponse, error) {
	attraction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.AttractionResponse{}, err
	}
	if attraction == nil {
		return dto.AttractionResponse{}, appError.NotFound(fmt.Sprintf("Attraction not found with ID: %d", id))
	}

	transports, err := s.transportRepo.GetAllByAttractionID(ctx, attraction.ID)
	if err != nil {
		return dto.AttractionResponse{}, err
	}

	return mapper.ToAttractionResponse(attraction, transports), nil
}

func (s *AttractionService) Create(ctx context.Context, req dto.AttractionRequest) (dto.AttractionResponse, error) {
	var createdAttraction *model.Attraction

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

		attraction := mapper.ToAttractionModel(req, sortOrder, imagePath)
		if err := txRepo.Create(ctx, attraction); err != nil {
			return err
		}

		createdAttraction = attraction

		return nil
	}); err != nil {
		return dto.AttractionResponse{}, err
	}

	return mapper.ToAttractionResponse(createdAttraction, nil), nil
}

func (s *AttractionService) Update(ctx context.Context, id uint, req dto.AttractionRequest) (dto.AttractionResponse, error) {
	var updatedAttraction *model.Attraction

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		attraction, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if attraction == nil {
			return appError.NotFound(fmt.Sprintf("Attraction not found with ID: %d", id))
		}

		oldOrder := attraction.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if *newOrder != oldOrder {
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

		imagepath := attraction.Image
		if req.Image != nil {
			if attraction.Image != nil {
				image.DeleteImage(*attraction.Image)
			}

			path, err := image.SaveImage("uploads/galleries", req.Image, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagepath = &path
		}

		mapper.UpdateAttractionModel(attraction, req, *newOrder, imagepath)
		if err := txRepo.Update(ctx, attraction); err != nil {
			return err
		}

		updatedAttraction = attraction

		return nil
	}); err != nil {
		return dto.AttractionResponse{}, err
	}

	return mapper.ToAttractionResponse(updatedAttraction, nil), nil
}

func (s *AttractionService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txtransportRepo := s.transportRepo.WithTx(tx)

		attraction, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if attraction == nil {
			return appError.NotFound(fmt.Sprintf("Attraction not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, attraction.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		if err := txtransportRepo.DeleteByAttractionID(ctx, id); err != nil {
			return err
		}

		return txRepo.Delete(ctx, attraction)
	})
}
