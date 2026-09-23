package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type ApplicationHandler struct {
	service *service.ApplicationService
}

func NewApplicationHandler(service *service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	ctx := c.Request.Context()

	application, err := h.service.Get(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Application successfully fetched!", application)
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req dto.ApplicationRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	application, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Application successfully created!", application)
}

func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	var req dto.ApplicationRequest

	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	application, err := h.service.Update(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Application successfully updated!", application)
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	ctx := c.Request.Context()

	if err := h.service.Delete(ctx); err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "Application successfully deleted!", nil)
}
