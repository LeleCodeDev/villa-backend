package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToLogoResponse(logo *model.Logo) dto.LogoResponse {
	return dto.LogoResponse{
		ID:        logo.ID,
		Text:      logo.Text,
		CreatedAt: logo.CreatedAt,
		UpdatedAt: logo.UpdatedAt,
	}
}

func ToLogoModel(req dto.LogoRequest) *model.Logo {
	return &model.Logo{
		Text: req.Text,
	}
}

func UpdateLogo(logo *model.Logo, req dto.LogoRequest) {
	logo.Text = req.Text
}
