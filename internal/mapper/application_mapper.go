// Package mapper
package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/media"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToApplicationResponse(application *model.Application) dto.ApplicationResponse {
	return dto.ApplicationResponse{
		ID:            application.ID,
		Title:         application.Title,
		PhoneNumber:   application.PhoneNumber,
		Logo:          media.BuildURLPtr(application.Logo),
		Music:         media.BuildURLPtr(application.Music),
		WeddingMusic:  media.BuildURLPtr(application.WeddingMusic),
		Email:         application.Email,
		InstagramLink: application.InstagramLink,
		TikTokLink:    application.TikTokLink,
		CreatedAt:     application.CreatedAt,
		UpdatedAt:     application.UpdatedAt,
	}
}

func ToApplicationModel(req dto.ApplicationRequest, filename *string, musicFileName *string, weddingMusicFileName *string) *model.Application {
	return &model.Application{
		Title:         req.Title,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		InstagramLink: req.InstagramLink,
		TikTokLink:    req.TikTokLink,
		Logo:          filename,
		Music:         musicFileName,
		WeddingMusic:  weddingMusicFileName,
	}
}

func UpdateApplicationModel(application *model.Application, req dto.ApplicationRequest, filename *string, musicFileName *string, weddingMusicFileName *string) {
	application.Title = req.Title
	application.PhoneNumber = req.PhoneNumber
	application.Logo = filename
	application.Music = musicFileName
	application.WeddingMusic = weddingMusicFileName
	application.Email = req.Email
	application.InstagramLink = req.InstagramLink
	application.TikTokLink = req.TikTokLink
}
