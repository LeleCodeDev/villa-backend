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

type TransportService struct {
	txManager      *repository.TxManager
	repo           *repository.TransportRepository
	attractionRepo *repository.AttractionRepository
}

func NewTransportService(
	txManager *repository.TxManager,
	repo *repository.TransportRepository,
	attractionRepo *repository.AttractionRepository,
) *TransportService {
	return &TransportService{
		txManager:      txManager,
		repo:           repo,
		attractionRepo: attractionRepo,
	}
}

func (s *TransportService) GetAll(ctx context.Context, query dto.TransportQuery) ([]dto.TransportResponse, int64, error) {
	transports, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.TransportResponse, 0, len(transports))
	for _, transport := range transports {
		responses = append(responses, mapper.ToTransportResponse(&transport))
	}

	return responses, total, nil
}

func (s *TransportService) GetAllByAttractionID(ctx context.Context, attractionID uint, query dto.TransportQuery) ([]dto.TransportResponse, int64, error) {
	transports, total, err := s.repo.GetAllByAttractionIDWithQuery(ctx, attractionID, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.TransportResponse, 0, len(transports))
	for _, transport := range transports {
		responses = append(responses, mapper.ToTransportResponse(&transport))
	}

	return responses, total, nil
}

func (s *TransportService) GetByID(ctx context.Context, id uint) (dto.TransportResponse, error) {
	transport, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.TransportResponse{}, err
	}
	if transport == nil {
		return dto.TransportResponse{}, appError.NotFound(fmt.Sprintf("Transport not found with ID: %d", id))
	}

	return mapper.ToTransportResponse(transport), nil
}

func (s *TransportService) Create(ctx context.Context, req dto.TransportRequest) (dto.TransportResponse, error) {
	var createdTransport *model.Transport

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txAttractionRepo := s.attractionRepo.WithTx(tx)

		attraction, err := txAttractionRepo.GetByID(ctx, req.AttractionID)
		if err != nil {
			return err
		}
		if attraction == nil {
			return appError.NotFound(fmt.Sprintf("Attraction not found with ID : %d", req.AttractionID))
		}

		var sortOrder int
		if req.SortOrder == nil {
			maxOrder, err := txRepo.GetMaxSortOrder(ctx, attraction.ID)
			if err != nil {
				return err
			}

			sortOrder = maxOrder + 1
		} else {
			sortOrder = *req.SortOrder
			if err := txRepo.ShiftSortOrder(ctx, attraction.ID, sortOrder); err != nil {
				return err
			}
		}

		transport := mapper.ToTransportModel(req, *attraction, sortOrder)
		if err := txRepo.Create(ctx, transport); err != nil {
			return err
		}

		createdTransport = transport

		return nil
	}); err != nil {
		return dto.TransportResponse{}, err
	}

	return mapper.ToTransportResponse(createdTransport), nil
}

func (s *TransportService) Update(ctx context.Context, id uint, req dto.TransportRequest) (dto.TransportResponse, error) {
	var updatedTransport *model.Transport

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txAttractionRepo := s.attractionRepo.WithTx(tx)

		transport, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if transport == nil {
			return appError.NotFound(fmt.Sprintf("Transport not found with ID: %d", id))
		}

		attraction := &transport.Attraction
		if transport.AttractionID != req.AttractionID {
			attraction, err = txAttractionRepo.GetByID(ctx, req.AttractionID)
			if err != nil {
				return err
			}
			if attraction == nil {
				return appError.NotFound(fmt.Sprintf("Attraction not found with ID : %d", req.AttractionID))
			}
		}

		oldOrder := transport.SortOrder
		newOrder := req.SortOrder

		if newOrder == nil {
			newOrder = &oldOrder
		} else if *newOrder != oldOrder {
			if *newOrder < oldOrder {
				if err := txRepo.UpdateSortOrderRange(ctx, attraction.ID, *newOrder, oldOrder-1, +1); err != nil {
					return err
				}
			} else {
				if err := txRepo.UpdateSortOrderRange(ctx, attraction.ID, oldOrder+1, *newOrder, -1); err != nil {
					return err
				}
			}
		}

		mapper.UpdateTransportModel(transport, *attraction, req, *newOrder)
		if err := txRepo.Update(ctx, transport); err != nil {
			return err
		}

		updatedTransport = transport

		return nil
	}); err != nil {
		return dto.TransportResponse{}, err
	}

	return mapper.ToTransportResponse(updatedTransport), nil
}

func (s *TransportService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		transport, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if transport == nil {
			return appError.NotFound(fmt.Sprintf("Transport not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx, transport.AttractionID)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, transport.AttractionID, transport.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, transport)
	})
}
