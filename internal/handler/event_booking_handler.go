package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/pagination"
	"github.com/lelecodedev/villa-backend/internal/response"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
)

type EventBookingHandler struct {
	service *service.EventBookingService
}

func NewEventBookingHandler(service *service.EventBookingService) *EventBookingHandler {
	return &EventBookingHandler{service: service}
}

func (h *EventBookingHandler) GetAllEventBookings(c *gin.Context) {
	var query dto.EventBookingQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	query.SetDefault()

	eventBookings, total, err := h.service.GetAll(ctx, query, nil)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All event bookings successfully fetched", eventBookings)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All event bookings successfully fetched", eventBookings, pagination)
}

func (h *EventBookingHandler) GetMyEventBookings(c *gin.Context) {
	var query dto.EventBookingQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	currentUser := c.MustGet("user").(model.User)
	query.SetDefault()

	eventBookings, total, err := h.service.GetAll(ctx, query, &currentUser.ID)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All event bookings successfully fetched", eventBookings)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All event bookings successfully fetched", eventBookings, pagination)
}

func (h *EventBookingHandler) GetEventBookingByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	eventBooking, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking successfully fetched", eventBooking)
}

func (h *EventBookingHandler) CreateEventBooking(c *gin.Context) {
	var req dto.EventBookingCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	currentUser := c.MustGet("user").(model.User)

	eventBooking, err := h.service.Create(ctx, currentUser, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Event booking successfully created!", eventBooking)
}

func (h *EventBookingHandler) CanceledEventBooking(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()
	currentUser := c.MustGet("user").(model.User)

	eventBooking, err := h.service.Cancel(ctx, currentUser, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking successfully canceled!", eventBooking)
}

func (h *EventBookingHandler) CheckInEventBooking(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()
	currentUser := c.MustGet("user").(model.User)

	eventBooking, err := h.service.CheckIn(ctx, currentUser, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking successfully check in!", eventBooking)
}

func (h *EventBookingHandler) CompletedEventBooking(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()
	currentUser := c.MustGet("user").(model.User)

	eventBooking, err := h.service.Complete(ctx, currentUser, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking successfully completed!", eventBooking)
}
