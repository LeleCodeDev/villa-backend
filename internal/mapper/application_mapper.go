// Package mapper
package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToApplicationResponse(application *model.Application) dto.ApplicationResponse {
	return dto.ApplicationResponse{
		ID:          application.ID,
		Title:       application.Title,
		PhoneNumber: application.PhoneNumber,
		Logo:        image.BuildURLPtr(application.Logo),
		CreatedAt:   application.CreatedAt,
		UpdatedAt:   application.UpdatedAt,
	}
}

func ToApplicationModel(req dto.ApplicationRequest, filename *string) *model.Application {
	return &model.Application{
		Title:       req.Title,
		PhoneNumber: req.PhoneNumber,
		Logo:        filename,
	}
}

func UpdateApplicationModel(application *model.Application, req dto.ApplicationRequest, filename *string) {
	application.Title = req.Title
	application.PhoneNumber = req.PhoneNumber
	application.Logo = filename
}
