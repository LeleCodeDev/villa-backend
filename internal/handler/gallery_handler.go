package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type GalleryHandler struct {
	service *service.GalleryService
}

func NewGalleryHandler(service *service.GalleryService) *GalleryHandler {
	return &GalleryHandler{service: service}
}

func (h *GalleryHandler) GetAllGalleries(c *gin.Context) {
	ctx := c.Request.Context()

	galleries, err := h.service.GetAll(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "All galleries successfully fetched!", galleries)
}

func (h *GalleryHandler) GetGalleryByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	gallery, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Gallery successfully fetched!", gallery)
}

func (h *GalleryHandler) CreateGallery(c *gin.Context) {
	var req dto.GalleryRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	gallery, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Gallery successfully created!", gallery)
}

func (h *GalleryHandler) UpdateGallery(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.GalleryRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	gallery, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Gallery successfully updated!", gallery)
}

func (h *GalleryHandler) DeleteGallery(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	if err := h.service.Delete(ctx, id); err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "Gallery successfully deleted!", nil)
}
