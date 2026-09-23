package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToVillaPackageListResponse(villaPackageList *model.VillaPackageList) dto.VillaPackageListResponse {
	return dto.VillaPackageListResponse{
		ID:        villaPackageList.ID,
		Text:      villaPackageList.Text,
		SortOrder: villaPackageList.SortOrder,
		CreatedAt: villaPackageList.CreatedAt,
		UpdatedAt: villaPackageList.UpdatedAt,
	}
}

func ToVillaPackageListModel(req dto.VillaPackageListRequest, sortOrder int) *model.VillaPackageList {
	return &model.VillaPackageList{
		Text:      req.Text,
		SortOrder: sortOrder,
	}
}

func UpdateVillaPackageListModel(villaPackageList *model.VillaPackageList, req dto.VillaPackageListRequest, sortOrder int) {
	villaPackageList.Text = req.Text
	villaPackageList.SortOrder = sortOrder
}
