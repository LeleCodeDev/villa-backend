package service

import "github.com/lelecodedev/villa-backend/internal/repository"

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
