package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToGalleryResponse(gallery *model.Gallery) dto.GalleryResponse {
	var imageURL *string
	if gallery.Image != nil {
		url := config.Env.BaseURL + "/api/" + *gallery.Image
		imageURL = &url
	}

	return dto.GalleryResponse{
		ID:          gallery.ID,
		Title:       gallery.Title,
		Description: gallery.Description,
		Image:       imageURL,
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
