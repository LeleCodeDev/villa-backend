package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/service"
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
