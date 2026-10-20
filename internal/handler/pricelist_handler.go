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

type PricelistHandler struct {
	service *service.PricelistService
}

func NewPricelistHandler(service *service.PricelistService) *PricelistHandler {
	return &PricelistHandler{service: service}
}

func (h *PricelistHandler) GetAllPricelists(c *gin.Context) {
	var query dto.PricelistQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleServiceError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	pricelists, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All pricelists successfully fetched!", pricelists)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All pricelists successfully fetched!", pricelists, pagination)
}

func (h *PricelistHandler) GetPricelistByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	pricelist, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Pricelist successfully fetched!", pricelist)
}

func (h *PricelistHandler) CreatePricelist(c *gin.Context) {
	var req dto.PricelistRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	pricelist, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Pricelist successfully created!", pricelist)
}

func (h *PricelistHandler) UpdatePricelist(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.PricelistRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	pricelist, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Pricelist successfully updated!", pricelist)
}

func (h *PricelistHandler) DeletePricelist(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Pricelist successfully deleted!", nil)
}
