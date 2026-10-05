package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToLogoResponse(logo *model.Logo) dto.LogoResponse {
	return dto.LogoResponse{
		ID:        logo.ID,
		Image:     image.BuildURL(logo.Image),
		CreatedAt: logo.CreatedAt,
		UpdatedAt: logo.UpdatedAt,
	}
}

func ToLogoModel(image string) *model.Logo {
	return &model.Logo{
		Image: image,
	}
}

func UpdateLogoModel(logo *model.Logo, image string) {
	logo.Image = image
}
