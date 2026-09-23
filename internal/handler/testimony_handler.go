package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type TestimonyHandler struct {
	service *service.TestimonyService
}

func NewTestimonyHandler(service *service.TestimonyService) *TestimonyHandler {
	return &TestimonyHandler{
		service: service,
	}
}

func (h *TestimonyHandler) GetAllTestimonies(c *gin.Context) {
	ctx := c.Request.Context()

	testimonies, err := h.service.GetAll(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
	}

	response.Success(c, http.StatusOK, "All testimonies successfully fetched!", testimonies)
}

func (h *TestimonyHandler) GetTestimonyByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
	}

	ctx := c.Request.Context()

	testimony, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Testimony successfully fetched!", testimony)
}

func (h *TestimonyHandler) CreateTestimony(c *gin.Context) {
	var req dto.TestimonyRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	testimony, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Testimony successfully created!", testimony)
}

func (h *TestimonyHandler) UpdateTestimony(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.TestimonyRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	testimony, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Testimony successfully updated!", testimony)
}

func (h *TestimonyHandler) DeleteTestimony(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Testimony successfully deleted!", nil)
}
