package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
)

type PaymentHandler struct{ service *service.PaymentService }

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Notification(c *gin.Context) {
	var req dto.MidtransNotification
	if err := c.ShouldBind(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.service.HandleNotification(c.Request.Context(), req)
	switch {
	case err == nil:
		c.Status(http.StatusOK)
	case errors.Is(err, service.ErrInvalidSignature):
		c.Status(http.StatusForbidden)
	case errors.Is(err, service.ErrBookingNotFound), errors.Is(err, service.ErrAmountMismatch):
		c.Status(http.StatusBadRequest)
	default:
		log.Println("payment notification:", err)
		c.Status(http.StatusInternalServerError)
	}
}
