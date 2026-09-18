package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type FaqHandler struct {
	service *service.FaqService
}

func NewFaqHandler(service *service.FaqService) *FaqHandler {
	return &FaqHandler{service: service}
}

func (h *FaqHandler) GetAllFaqs(c *gin.Context) {
	ctx := c.Request.Context()

	faqs, err := h.service.GetAll(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "All faqs successfully fetched!", faqs)
}
