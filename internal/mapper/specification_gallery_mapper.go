package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToSpecificationGalleryResponse(specificationGallery *model.SpecificationGallery) dto.SpecificationGalleryResponse {
	return dto.SpecificationGalleryResponse{
		ID:        specificationGallery.ID,
		Image1:    image.BuildURLPtr(specificationGallery.Image1),
		Image2:    image.BuildURLPtr(specificationGallery.Image2),
		Image3:    image.BuildURLPtr(specificationGallery.Image3),
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
