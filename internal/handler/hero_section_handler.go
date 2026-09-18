package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type HeroSectionHandler struct {
	service *service.HeroSectionService
}

func NewHeroSectionHandler(service *service.HeroSectionService) *HeroSectionHandler {
	return &HeroSectionHandler{service: service}
}

func (h *HeroSectionHandler) GetHeroSecton(c *gin.Context) {
	ctx := c.Request.Context()

	heroSection, err := h.service.Get(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Hero section successfully fetched!", heroSection)
}

func (h *HeroSectionHandler) CreateHeroSection(c *gin.Context) {
	var req dto.HeroSectionRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	heroSection, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Hero section successfully created!", heroSection)
}

func (h *HeroSectionHandler) UpdateHeroSection(c *gin.Context) {
	var req dto.HeroSectionRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	heroSection, err := h.service.Update(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Hero section successfully updated!", heroSection)
}
