package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/pagination"
	"github.com/lelecodedev/villa-backend/internal/response"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	var query dto.BookingQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	query.SetDefault()

	bookings, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All booking successfully fetched!", bookings)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All booking successfully fetched!", bookings, pagination)
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	booking, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Booking successfully fetched!", booking)
}
