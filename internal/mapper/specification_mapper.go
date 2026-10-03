package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToSpecificationResponse(specification *model.Specification) dto.SpecificationResponse {
	return dto.SpecificationResponse{
		ID:        specification.ID,
		Text:      specification.Text,
		Logo:      ToLogoResponse(&specification.Logo),
		SortOrder: specification.SortOrder,
		CreatedAt: specification.CreatedAt,
		UpdatedAt: specification.UpdatedAt,
	}
}

func ToSpecificationModel(req dto.SpecificationRequest, sortOrder int, logo *model.Logo) *model.Specification {
	return &model.Specification{
		Text:      req.Text,
		LogoID:    logo.ID,
		Logo:      *logo,
		SortOrder: sortOrder,
	}
}

func UpdateSpecificationModel(specification *model.Specification, req dto.SpecificationRequest, logo *model.Logo, sortOrder int) {
	specification.Text = req.Text
	specification.LogoID = req.LogoID
	specification.Logo = *logo
	specification.SortOrder = sortOrder
}
