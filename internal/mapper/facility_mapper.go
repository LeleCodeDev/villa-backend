package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToFacilityResponse(facility *model.Facility) dto.FacilityResponse {
	return dto.FacilityResponse{
		ID:          facility.ID,
		Title:       facility.Title,
		Description: facility.Description,
		Logo:        facility.Logo,
		SortOrder:   facility.SortOrder,
		CreatedAt:   facility.CreatedAt,
		UpdatedAt:   facility.UpdatedAt,
	}
}

func ToFacilityModel(req dto.FacilityRequest, sortOrder int) *model.Facility {
	return &model.Facility{
		Title:       req.Title,
		Description: req.Description,
		Logo:        req.Logo,
		SortOrder:   sortOrder,
	}
}

func UpdateFacilityModel(facility *model.Facility, req dto.FacilityRequest, sortOrder int) {
	facility.Title = req.Title
	facility.Description = req.Description
	facility.Logo = req.Logo
	facility.SortOrder = sortOrder
}
