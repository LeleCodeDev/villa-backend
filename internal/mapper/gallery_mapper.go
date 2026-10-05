package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToGalleryResponse(gallery *model.Gallery) dto.GalleryResponse {
	return dto.GalleryResponse{
		ID:          gallery.ID,
		Title:       gallery.Title,
		Description: gallery.Description,
		Image:       image.BuildURLPtr(gallery.Image),
		SortOrder:   gallery.SortOrder,
		CreatedAt:   gallery.CreatedAt,
		UpdatedAt:   gallery.UpdatedAt,
	}
}

func ToGalleryModel(req dto.GalleryRequest, filename *string, order int) *model.Gallery {
	return &model.Gallery{
		Title:       req.Title,
		Description: req.Description,
		Image:       filename,
		SortOrder:   order,
	}
}

func UpdateGalleryModel(gallery *model.Gallery, req dto.GalleryRequest, filename *string, order int) {
	gallery.Title = req.Title
	gallery.Description = req.Description
	gallery.Image = filename
	gallery.SortOrder = order
}
