package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
	"github.com/lelecodedev/villa-backend/internal/pagination"
	"github.com/lelecodedev/villa-backend/internal/response"
)

type FacilityHandler struct {
	service *service.FacilityService
}

func NewFacilityHandler(service *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{service: service}
}

func (h *FacilityHandler) GetAllFacilities(c *gin.Context) {
	var query dto.FacilityQuery

	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	facilities, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All facilities successfully fetched!", facilities)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All facilities successfully fetched!", facilities, pagination)
}

func (h *FacilityHandler) GetFacilityByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	facility, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facility successfully fetched!", facility)
}

func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req dto.FacilityRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	facility, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Facility successfully created!", facility)
}

func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.FacilityRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	facility, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Facility successfully updated!", facility)
}

func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Facility successfully deleted!", nil)
}
