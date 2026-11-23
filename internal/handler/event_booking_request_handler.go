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

type EventBookingRequestHandler struct {
	service *service.EventBookingRequestService
}

func NewEventBookingRequestHandler(service *service.EventBookingRequestService) *EventBookingRequestHandler {
	return &EventBookingRequestHandler{service: service}
}

func (h *EventBookingRequestHandler) GetAllEventBookingRequests(c *gin.Context) {
	var query dto.EventBookingRequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()
	query.SetDefault()

	requests, total, err := h.service.GetAll(ctx, query, nil)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All event booking requests successfully fetched!", requests)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All event booking requests successfully fetched!", requests, pagination)
}

func (h *EventBookingRequestHandler) GetMyEventBookingRequests(c *gin.Context) {
	var query dto.EventBookingRequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()
	query.SetDefault()

	requests, total, err := h.service.GetAll(ctx, query, &currentUser.ID)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All event booking requests successfully fetched!", requests)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All event booking requests successfully fetched!", requests, pagination)
}

func (h *EventBookingRequestHandler) GetEventBookingRequestByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()

	request, err := h.service.GetByID(ctx, id, currentUser)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking request successfully fetched!", request)
}

func (h *EventBookingRequestHandler) CreateEventBookingRequest(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()

	request, err := h.service.Create(ctx, req, currentUser)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Event booking request successfully created!", request)
}

func (h *EventBookingRequestHandler) ApproveEventBookingRequest(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.ApproveEventRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	request, err := h.service.Approve(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking request successfully approved!", request)
}

func (h *EventBookingRequestHandler) RejectEventBookingRequest(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.RejectEventRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	request, err := h.service.Reject(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking request successfully rejected!", request)
}

func (h *EventBookingRequestHandler) PayEventBookingRequest(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()

	request, err := h.service.Pay(ctx, currentUser, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking request successfully rejected!", request)
}

func (h *EventBookingRequestHandler) CancelEventBookingRequest(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()

	request, err := h.service.Cancel(ctx, currentUser, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Event booking request successfully rejected!", request)
}
