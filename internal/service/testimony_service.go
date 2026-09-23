package service

import (
	"context"
	"fmt"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"gorm.io/gorm"
)

type TestimonyService struct {
	txManager *repository.TxManager
	repo      *repository.TestimonyRepository
}

func NewTestimonyService(
	txManager *repository.TxManager,
	repo *repository.TestimonyRepository,
) *TestimonyService {
	return &TestimonyService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *TestimonyService) GetAll(ctx context.Context) ([]dto.TestimonyResponse, error) {
	testimonies, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TestimonyResponse, 0, len(testimonies))
	for _, testimony := range testimonies {
		responses = append(responses, mapper.ToTestimonyResponse(&testimony))
	}

	return responses, nil
}

func (s *TestimonyService) GetByID(ctx context.Context, id uint) (dto.TestimonyResponse, error) {
	testimony, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.TestimonyResponse{}, err
	}
	if testimony == nil {
		return dto.TestimonyResponse{}, errors.NotFound(fmt.Sprintf("Testimony not found with ID : %d", id))
	}

	return mapper.ToTestimonyResponse(testimony), nil
}

func (s *TestimonyService) Create(ctx context.Context, req dto.TestimonyRequest) (dto.TestimonyResponse, error) {
	var createdTestimony *model.Testimony

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

		testimony := mapper.ToTestimonyModel(req, sortOrder)
		if err := txRepo.Create(ctx, testimony); err != nil {
			return err
		}

		createdTestimony = testimony

		return nil
	}); err != nil {
		return dto.TestimonyResponse{}, err
	}

	return mapper.ToTestimonyResponse(createdTestimony), nil
}

func (s *TestimonyService) Update(ctx context.Context, req dto.TestimonyRequest, id uint) (dto.TestimonyResponse, error) {
	var updatedTestimony *model.Testimony

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		testimony, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if testimony == nil {
			return errors.NotFound(fmt.Sprintf("Testimony not found with ID: %d", id))
		}

		oldOrder := testimony.SortOrder
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

		mapper.UpdateTestimonyModel(testimony, req, *newOrder)
		if err := txRepo.Update(ctx, testimony); err != nil {
			return err
		}

		updatedTestimony = testimony

		return nil
	}); err != nil {
		return dto.TestimonyResponse{}, err
	}

	return mapper.ToTestimonyResponse(updatedTestimony), nil
}

func (s *TestimonyService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		testimony, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if testimony == nil {
			return errors.NotFound(fmt.Sprintf("Testimony not found with ID: %d", id))
		}

		maxOrder, err := txRepo.GetMaxSortOrder(ctx)
		if err != nil {
			return err
		}
		if err := txRepo.UpdateSortOrderRange(ctx, testimony.SortOrder+1, maxOrder, -1); err != nil {
			return err
		}

		return txRepo.Delete(ctx, testimony)
	})
}
