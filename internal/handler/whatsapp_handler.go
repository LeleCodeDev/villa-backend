package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/response"
	"github.com/lelecodedev/villa-backend/internal/service"
)

type WhatsappHandler struct {
	service *service.WhatsappService
}

func NewWhatsappHandler(service *service.WhatsappService) *WhatsappHandler {
	return &WhatsappHandler{service: service}
}

func (h *WhatsappHandler) SendText(c *gin.Context) {
	var req dto.SendWhatsAppRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	if err := h.service.Send(ctx, req.Phone, req.Message); err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "Message successfully send!", nil)
}
