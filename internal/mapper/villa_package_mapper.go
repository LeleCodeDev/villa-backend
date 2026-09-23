package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToVillaPackageResponse(villaPackage *model.VillaPackage, villaPackageResponses []dto.VillaPackageListResponse) dto.VillaPackageResponse {
	return dto.VillaPackageResponse{
		ID:          villaPackage.ID,
		Title:       villaPackage.Title,
		Subtitle:    villaPackage.Subtitle,
		MaxCapacity: villaPackage.MaxCapacity,
		Price:       villaPackage.Price,
		Lists:       villaPackageResponses,
		SortOrder:   villaPackage.SortOrder,
		CreatedAt:   villaPackage.CreatedAt,
		UpdatedAt:   villaPackage.UpdatedAt,
	}
}

func ToVillaPackageModel(req dto.VillaPackageRequest, sortOrder int) *model.VillaPackage {
	return &model.VillaPackage{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		MaxCapacity: req.MaxCapacity,
		Price:       req.Price,
		SortOrder:   sortOrder,
	}
}

func UpdateVillaPackageModel(villaPackage *model.VillaPackage, req dto.VillaPackageRequest, sortOrder int) {
	villaPackage.Title = req.Title
	villaPackage.Subtitle = req.Subtitle
	villaPackage.Price = req.Price
	villaPackage.MaxCapacity = req.MaxCapacity
	villaPackage.SortOrder = sortOrder
}
