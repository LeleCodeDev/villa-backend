package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToRentalOptionResponse(rentalOption *model.RentalOption) dto.RentalOptionResponse {
	return dto.RentalOptionResponse{
		ID:          rentalOption.ID,
		Title:       rentalOption.Title,
		Subtitle:    rentalOption.Subtitle,
		Description: rentalOption.Description,
		RoomCount:   rentalOption.RoomCount,
		MaxCapacity: rentalOption.MaxCapacity,
		SortOrder:   rentalOption.SortOrder,
		Image:       image.BuildURLPtr(rentalOption.Image),
		CreatedAt:   rentalOption.CreatedAt,
		UpdatedAt:   rentalOption.UpdatedAt,
	}
}

func ToRentalOptionModel(req dto.RentalOptionRequest, sortOrder int, imagePath *string) *model.RentalOption {
	return &model.RentalOption{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		RoomCount:   req.RoomCount,
		MaxCapacity: req.MaxCapacity,
		SortOrder:   sortOrder,
		Image:       imagePath,
	}
}

func UpdateRentalOptionModel(rentalOption *model.RentalOption, req dto.RentalOptionRequest, sortOrder int, imagePath *string) {
	rentalOption.Title = req.Title
	rentalOption.Subtitle = req.Subtitle
	rentalOption.Description = req.Description
	rentalOption.RoomCount = req.RoomCount
	rentalOption.MaxCapacity = req.MaxCapacity
	rentalOption.SortOrder = sortOrder
	rentalOption.Image = imagePath
}
