// Package mapper
package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/image"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToApplicationResponse(application *model.Application) dto.ApplicationResponse {
	return dto.ApplicationResponse{
		ID:            application.ID,
		Title:         application.Title,
		PhoneNumber:   application.PhoneNumber,
		Logo:          image.BuildURLPtr(application.Logo),
		Email:         application.Email,
		InstagramLink: application.InstagramLink,
		TikTokLink:    application.TikTokLink,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     application.UpdatedAt,
	}
}

func ToApplicationModel(req dto.ApplicationRequest, filename *string) *model.Application {
	return &model.Application{
		Title:         req.Title,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		InstagramLink: req.InstagramLink,
		TikTokLink:    req.TikTokLink,
		Logo:          filename,
	}
}

func UpdateApplicationModel(application *model.Application, req dto.ApplicationRequest, filename *string) {
	application.Title = req.Title
	application.PhoneNumber = req.PhoneNumber
	application.Logo = filename
	application.Email = req.Email
	application.InstagramLink = req.InstagramLink
	application.TikTokLink = req.TikTokLink
}
