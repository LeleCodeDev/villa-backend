package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToVillaPackageResponse(villaPackage *model.VillaPackage, lists []model.VillaPackageList) dto.VillaPackageResponse {
	var listsResponse []dto.VillaPackageListResponse
	for i := range lists {
		listsResponse = append(listsResponse, ToVillaPackageListResponse(&lists[i]))
	}

	return dto.VillaPackageResponse{
		ID:               villaPackage.ID,
		Title:            villaPackage.Title,
		Subtitle:         villaPackage.Subtitle,
		Type:             villaPackage.Type,
		Price:            villaPackage.Price,
		RequiresApproval: villaPackage.RequiresApproval,
		MaxCapacity:      villaPackage.MaxCapacity,
		IsWedding:        villaPackage.IsWedding,
		Lists:            listsResponse,
		SortOrder:        villaPackage.SortOrder,
		CreatedAt:        villaPackage.CreatedAt,
		UpdatedAt:        villaPackage.UpdatedAt,
	}
}

func ToVillaPackageModel(req dto.VillaPackageRequest, sortOrder int, custom bool) *model.VillaPackage {
	return &model.VillaPackage{
		Title:            req.Title,
		Subtitle:         req.Subtitle,
		Type:             req.Type,
		Price:            req.Price,
		MaxCapacity:      req.MaxCapacity,
		IsWedding:        req.IsWedding,
		RequiresApproval: *req.RequiresApproval,
		SortOrder:        sortOrder,
	}
}

func UpdateVillaPackageModel(villaPackage *model.VillaPackage, req dto.VillaPackageRequest, sortOrder int) {
	villaPackage.Title = req.Title
	villaPackage.Subtitle = req.Subtitle
	villaPackage.Price = req.Price
	villaPackage.Type = req.Type
	villaPackage.IsWedding = req.IsWedding
	villaPackage.SortOrder = sortOrder
}
