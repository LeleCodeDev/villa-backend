package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
	"github.com/lelecodedev/villa-backend/pkg/pagination"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type AddonListHandler struct {
	service *service.AddonListService
}

func NewAddonListHandler(service *service.AddonListService) *AddonListHandler {
	return &AddonListHandler{service: service}
}

func (h *AddonListHandler) GetAllAddonLists(c *gin.Context) {
	var query dto.AddonListQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	addonLists, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All addon lists successfully fetched!", addonLists)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All addon lists successfully fetched!", addonLists, pagination)
}

func (h *AddonListHandler) GetAddonListByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	addonList, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Addon list successfully fetched!", addonList)
}

func (h *AddonListHandler) CreateAddonList(c *gin.Context) {
	var req dto.AddonListRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	addonList, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Addon list successfully created!", addonList)
}

func (h *AddonListHandler) UpdateAddonList(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.AddonListRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	addonList, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Addon list successfully updated!", addonList)
}

func (h *AddonListHandler) DeleteAddonList(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Addon list successfully deleted!", nil)
}
