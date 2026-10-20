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

type ReasonHandler struct {
	service *service.ReasonService
}

func NewReasonHandler(service *service.ReasonService) *ReasonHandler {
	return &ReasonHandler{service: service}
}

func (h *ReasonHandler) GetAllReasons(c *gin.Context) {
	var query dto.ReasonQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	reasons, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {

		response.Success(c, http.StatusOK, "All reasons successfully fetched!", reasons)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All reasons successfully fetched!", reasons, pagination)
}

func (h *ReasonHandler) GetReasonByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	reason, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Reason successfully fetched!", reason)
}

func (h *ReasonHandler) CreateReason(c *gin.Context) {
	var req dto.ReasonRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	reason, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Reason successfully created!", reason)
}

func (h *ReasonHandler) UpdateReason(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.ReasonRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	reason, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Reason successfully updated!", reason)
}

func (h *ReasonHandler) DeleteReason(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Reason successfully updated!", nil)
}
