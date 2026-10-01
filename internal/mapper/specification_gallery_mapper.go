package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToSpecificationGalleryResponse(specificationGallery *model.SpecificationGallery) dto.SpecificationGalleryResponse {
	var imageURL1 *string
	var imageURL2 *string
	var imageURL3 *string
	if specificationGallery.Image1 != nil {
		url := config.Env.BaseURL + "/api/" + *specificationGallery.Image1
		imageURL1 = &url
	}
	if specificationGallery.Image2 != nil {
		url := config.Env.BaseURL + "/api/" + *specificationGallery.Image2
		imageURL2 = &url
	}
	if specificationGallery.Image3 != nil {
		url := config.Env.BaseURL + "/api/" + *specificationGallery.Image3
		imageURL3 = &url
	}

	return dto.SpecificationGalleryResponse{
		ID:        specificationGallery.ID,
		Image1:    imageURL1,
		Image2:    imageURL2,
		Image3:    imageURL3,
		CreatedAt: specificationGallery.CreatedAt,
		UpdatedAt: specificationGallery.UpdatedAt,
	}
}

func ToSpecificationGalleryModel(
	filename1 *string,
	filename2 *string,
	filename3 *string,
) *model.SpecificationGallery {
	return &model.SpecificationGallery{
		Image1: filename1,
		Image2: filename2,
		Image3: filename3,
	}
}

func UpdateSpecificationGalleryModel(
	specificationGallery *model.SpecificationGallery,
	filename1 *string,
	filename2 *string,
	filename3 *string,
) {
	specificationGallery.Image1 = filename1
	specificationGallery.Image2 = filename2
	specificationGallery.Image3 = filename3
}
