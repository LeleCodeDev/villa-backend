package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"gorm.io/gorm"
)

type HeroSectionService struct {
	txManager *repository.TxManager
	repo      *repository.HeroSectionRepository
}

func NewHeroSectionService(
	txManager *repository.TxManager,
	repo *repository.HeroSectionRepository,
) *HeroSectionService {
	return &HeroSectionService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *HeroSectionService) Get(ctx context.Context) (dto.HeroSectionResponse, error) {
	heroSection, err := s.repo.Get(ctx)
	if err != nil {
		return dto.HeroSectionResponse{}, err
	}
	if heroSection == nil {
		return dto.HeroSectionResponse{}, errors.NotFound("Hero section data not found!")
	}

	return mapper.ToHeroSectionResponse(heroSection), nil
}

func (s *HeroSectionService) Create(ctx context.Context, req dto.HeroSectionRequest) (dto.HeroSectionResponse, error) {
	var createdHeroSection *model.HeroSection

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			return errors.AlreadyExist("Hero section data already exist!")
		}

		heroSection := mapper.ToHeroSectionModel(req)
		if err := txRepo.Create(ctx, heroSection); err != nil {
			return err
		}

		createdHeroSection = heroSection

		return nil
	}); err != nil {
		return dto.HeroSectionResponse{}, err
	}

	return mapper.ToHeroSectionResponse(createdHeroSection), nil
}

func (s *HeroSectionService) Update(ctx context.Context, req dto.HeroSectionRequest) (dto.HeroSectionResponse, error) {
	var updatedHeroSection *model.HeroSection

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		heroSection, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if heroSection == nil {
			return errors.NotFound("Hero section data not found!")
		}

		mapper.UpdateHeroSectionModel(heroSection, req)
		if err := txRepo.Update(ctx, heroSection); err != nil {
			return err
		}

		updatedHeroSection = heroSection

		return nil
	}); err != nil {
		return dto.HeroSectionResponse{}, err
	}

	return mapper.ToHeroSectionResponse(updatedHeroSection), nil
}

func (s *HeroSectionService) Delete(ctx context.Context) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		heroSection, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if heroSection == nil {
			return errors.NotFound("Hero section data not found!")
		}

		return txRepo.Delete(ctx, heroSection)
	})
}
