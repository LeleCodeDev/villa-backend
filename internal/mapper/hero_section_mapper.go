package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToHeroSectionResponse(heroSection *model.HeroSection) dto.HeroSectionResponse {
	return dto.HeroSectionResponse{
		ID:        heroSection.ID,
		Title:     heroSection.Title,
		Subtitle:  heroSection.Subtitle,
		CreatedAt: heroSection.CreatedAt,
		UpdatedAt: heroSection.UpdatedAt,
	}
}

func ToHeroSectionModel(req dto.HeroSectionRequest) *model.HeroSection {
	return &model.HeroSection{
		Title:    req.Title,
		Subtitle: req.Subtitle,
	}
}

func UpdateHeroSectionModel(heroSection *model.HeroSection, req dto.HeroSectionRequest) {
	heroSection.Title = req.Title
	heroSection.Subtitle = req.Subtitle
}
