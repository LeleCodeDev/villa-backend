package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/response"
)

type SpecificationGalleryHandler struct {
	service *service.SpecificationGalleryService
}

func NewSpecificationGalleryHandler(service *service.SpecificationGalleryService) *SpecificationGalleryHandler {
	return &SpecificationGalleryHandler{service: service}
}

func (h *SpecificationGalleryHandler) GetSpecificationGallery(c *gin.Context) {
	ctx := c.Request.Context()

	specificationGallery, err := h.service.Get(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Specification gallery successfully fetched!", specificationGallery)
}

func (h *SpecificationGalleryHandler) CreateSpecificationGallery(c *gin.Context) {
	var req dto.SpecificationGalleryRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	specificationGallery, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Specification gallery successfully created!", specificationGallery)
}

func (h *SpecificationGalleryHandler) UpdateSpecificationGallery(c *gin.Context) {
	var req dto.SpecificationGalleryRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	specificationGallery, err := h.service.Update(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Specification gallery successfully updated!", specificationGallery)
}

func (h *SpecificationGalleryHandler) DeleteSpecificationGallery(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.service.Delete(ctx); err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "Specification gallery successfully deleted!", nil)
}
