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

type SeasonPeriodHandler struct {
	service *service.SeasonPeriodService
}

func NewSeasonPeriodHandler(service *service.SeasonPeriodService) *SeasonPeriodHandler {
	return &SeasonPeriodHandler{service: service}
}

func (h *SeasonPeriodHandler) GetAllSeasonPeriods(c *gin.Context) {
	var query dto.SeasonPeriodQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	seasonPeriods, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All season period successfully fetched!", seasonPeriods)
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All season period successfully fetched!", seasonPeriods, pagination)
}

func (h *SeasonPeriodHandler) GetSeasonPeriodByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	seasonPeriod, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Season period successfully fetched!", seasonPeriod)
}

func (h *SeasonPeriodHandler) CreateSeasonPeriod(c *gin.Context) {
	var req dto.SeasonPeriodRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	seasonPeriod, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Season period successfully created!", seasonPeriod)
}

func (h *SeasonPeriodHandler) UpdateSeasonPeriod(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.SeasonPeriodRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	seasonPeriod, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Season period successfully updated!", seasonPeriod)
}
