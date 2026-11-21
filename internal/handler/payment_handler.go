package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
)

type PaymentHandler struct{ service *service.PaymentService }

func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Notification(c *gin.Context) {
	var req dto.MidtransCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()

	err := h.service.HandleNotificationCallback(ctx, req)
	if err != nil {
		if serviceError, ok := err.(*appError.ServiceError); ok {
			log.Println(serviceError.Message)
			c.Status(serviceError.StatusCode)
			return
		}

		log.Println("Payment notification:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
