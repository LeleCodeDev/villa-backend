package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToSpecificationResponse(specification *model.Specification) dto.SpecificationResponse {
	return dto.SpecificationResponse{
		ID:        specification.ID,
		Text:      specification.Text,
		Logo:      specification.Logo,
		SortOrder: specification.SortOrder,
		CreatedAt: specification.CreatedAt,
		UpdatedAt: specification.UpdatedAt,
	}
}

func ToSpecificationModel(req dto.SpecificationRequest, sortOrder int) *model.Specification {
	return &model.Specification{
		Text:      req.Text,
		Logo:      req.Logo,
		SortOrder: sortOrder,
	}
}

func UpdateSpecificationModel(specification *model.Specification, req dto.SpecificationRequest, sortOrder int) {
	specification.Text = req.Text
	specification.Logo = req.Logo
	specification.SortOrder = sortOrder
}
