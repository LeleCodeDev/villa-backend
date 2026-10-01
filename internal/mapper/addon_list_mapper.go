package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToAddonListResponse(addonList *model.AddonList) dto.AddonListResponse {
	return dto.AddonListResponse{
		ID:        addonList.ID,
		Text:      addonList.Text,
		SortOrder: addonList.SortOrder,
		CreatedAt: addonList.CreatedAt,
		UpdatedAt: addonList.UpdatedAt,
	}
}

func ToAddonListModel(req dto.AddonListRequest, addon model.Addon, sortOrder int) *model.AddonList {
	return &model.AddonList{
		Text:      req.Text,
		Addon:     addon,
		AddonID:   addon.ID,
		SortOrder: sortOrder,
	}
}

func UpdateAddonListModel(addonList *model.AddonList, addon model.Addon, req dto.AddonListRequest, sortOrder int) {
	addonList.Text = req.Text
	addonList.AddonID = addon.ID
	addonList.Addon = addon
	addonList.SortOrder = sortOrder
}
