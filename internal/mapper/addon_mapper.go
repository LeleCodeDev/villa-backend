package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToAddonResponse(addon *model.Addon, listResponses []dto.AddonListResponse) dto.AddonResponse {
	return dto.AddonResponse{
		ID:          addon.ID,
		Title:       addon.Title,
		Subtitle:    addon.Subtitle,
		Description: addon.Description,
		Price:       addon.Price,
		Lists:       listResponses,
		SortOrder:   addon.SortOrder,
		CreatedAt:   addon.CreatedAt,
		UpdatedAt:   addon.UpdatedAt,
	}
}

func ToAddonModel(req dto.AddonRequest, sortOrder int) *model.Addon {
	return &model.Addon{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Price:       req.Price,
		SortOrder:   sortOrder,
	}
}

func UpdateAddonModel(addon *model.Addon, req dto.AddonRequest, sortOrder int) {
	addon.Title = req.Title
	addon.Subtitle = req.Subtitle
	addon.Price = req.Price
	addon.Description = req.Description
	addon.SortOrder = sortOrder
}
