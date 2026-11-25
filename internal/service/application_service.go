package service

import (
	"context"
	"log"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/audio"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/media"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
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
		return dto.ApplicationResponse{}, appError.NotFound("Application data not found!")
	}

	return mapper.ToApplicationResponse(application), nil
}

func (s *ApplicationService) Create(ctx context.Context, req dto.ApplicationRequest) (dto.ApplicationResponse, error) {
	var createdApplication *model.Application
	var imagePath *string
	var musicPath *string
	var weddingMusicPath *string

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist("Application data already exist!")
		}

		if req.Logo != nil {
			path, err := image.Save("uploads/application", req.Logo, image.DefaultOptions())
			if err != nil {
				return err
			}

			imagePath = &path
		}

		if req.Music != nil {
			path, err := audio.Save("uploads/application", req.Music, audio.DefaultOptions())
			if err != nil {
				return err
			}

			musicPath = &path
		}

		if req.WeddingMusic != nil {
			path, err := audio.Save("uploads/application", req.WeddingMusic, audio.DefaultOptions())
			if err != nil {
				return err
			}

			weddingMusicPath = &path
		}

		application := mapper.ToApplicationModel(req, imagePath, musicPath, weddingMusicPath)
		if err := txRepo.Create(ctx, application); err != nil {
			return err
		}

		createdApplication = application

		return nil
	}); err != nil {
		if imagePath != nil {
			if delErr := media.Delete(*imagePath); delErr != nil {
				log.Printf("failed to cleanup logo %s: %v", *imagePath, delErr)
			}
		}

		if musicPath != nil {
			if delErr := media.Delete(*musicPath); delErr != nil {
				log.Printf("failed to cleanup music %s: %v", *musicPath, delErr)
			}
		}

		if weddingMusicPath != nil {
			if delErr := media.Delete(*weddingMusicPath); delErr != nil {
				log.Printf("failed to cleanup wedding music %s: %v", *weddingMusicPath, delErr)
			}
		}

		return dto.ApplicationResponse{}, err
	}

	return mapper.ToApplicationResponse(createdApplication), nil
}

func (s *ApplicationService) Update(ctx context.Context, req dto.ApplicationRequest) (dto.ApplicationResponse, error) {
	var updatedApplication *model.Application
	var oldLogo *string
	var newLogo *string

	var oldMusic *string
	var newMusic *string

	var oldWeddingMusic *string
	var newWeddingMusic *string

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		application, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if application == nil {
			return appError.NotFound("Application data not found!")
		}

		logoPath := application.Logo
		if req.Logo != nil {
			path, err := image.Save("uploads/application", req.Logo, image.DefaultOptions())
			if err != nil {
				return err
			}
			newLogo = &path
			oldLogo = application.Logo
			logoPath = &path
		}

		musicPath := application.Music
		if req.Music != nil {
			path, err := image.Save("uploads/application", req.Music, image.DefaultOptions())
			if err != nil {
				return err
			}
			newMusic = &path
			oldMusic = application.Music
			musicPath = &path
		}

		weddingMusicPath := application.WeddingMusic
		if req.WeddingMusic != nil {
			path, err := image.Save("uploads/application", req.WeddingMusic, image.DefaultOptions())
			if err != nil {
				return err
			}
			newWeddingMusic = &path
			oldWeddingMusic = application.WeddingMusic
			weddingMusicPath = &path
		}

		mapper.UpdateApplicationModel(application, req, logoPath, musicPath, weddingMusicPath)
		if err := txRepo.Update(ctx, application); err != nil {
			return err
		}

		updatedApplication = application
		return nil
	}); err != nil {
		if newLogo != nil {
			if delErr := media.Delete(*newLogo); delErr != nil {
				log.Printf("failed to cleanup new logo %s: %v", *newLogo, delErr)
			}
		}

		if newMusic != nil {
			if delErr := media.Delete(*newMusic); delErr != nil {
				log.Printf("failed to cleanup new music %s: %v", *newMusic, delErr)
			}
		}

		if newWeddingMusic != nil {
			if delErr := media.Delete(*newWeddingMusic); delErr != nil {
				log.Printf("failed to cleanup new wedding music %s: %v", *newWeddingMusic, delErr)
			}
		}

		return dto.ApplicationResponse{}, err
	}

	if oldLogo != nil {
		if delErr := media.Delete(*oldLogo); delErr != nil {
			log.Printf("failed to delete old logo %s: %v", *oldLogo, delErr)
		}
	}

	if oldMusic != nil {
		if delErr := media.Delete(*oldMusic); delErr != nil {
			log.Printf("failed to delete old music %s: %v", *oldMusic, delErr)
		}
	}

	if oldWeddingMusic != nil {
		if delErr := media.Delete(*oldWeddingMusic); delErr != nil {
			log.Printf("failed to delete old wedding music %s: %v", *oldWeddingMusic, delErr)
		}
	}

	return mapper.ToApplicationResponse(updatedApplication), nil
}

func (s *ApplicationService) Delete(ctx context.Context) error {
	var imagePath *string
	var musicPath *string
	var weddingMusicPath *string

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		application, err := txRepo.Get(ctx)
		if err != nil {
			return err
		}
		if application == nil {
			return appError.NotFound("Application data not found!")
		}

		imagePath = application.Logo
		musicPath = application.Music
		weddingMusicPath = application.WeddingMusic

		return txRepo.Delete(ctx, application)
	}); err != nil {
		return err
	}

	if imagePath != nil {
		if delErr := media.Delete(*imagePath); delErr != nil {
			log.Printf("failed to delete logo %s: %v", *imagePath, delErr)
		}
	}

	if musicPath != nil {
		if delErr := media.Delete(*musicPath); delErr != nil {
			log.Printf("failed to delete music %s: %v", *musicPath, delErr)
		}
	}

	if weddingMusicPath != nil {
		if delErr := media.Delete(*weddingMusicPath); delErr != nil {
			log.Printf("failed to delete wedding music %s: %v", *weddingMusicPath, delErr)
		}
	}

	return nil
}
