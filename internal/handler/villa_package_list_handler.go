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

type VillaPackageListHandler struct {
	service *service.VillaPackageListService
}

func NewVillaPackageListHandler(service *service.VillaPackageListService) *VillaPackageListHandler {
	return &VillaPackageListHandler{service: service}
}

func (h *VillaPackageListHandler) GetAllVillaPackageLists(c *gin.Context) {
	var query dto.VillaPackageListQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	villaPackageLists, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All villa package lists successfully fetched!", villaPackageLists)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All villa package lists successfully fetched!", villaPackageLists, pagination)
}

func (h *VillaPackageListHandler) GetVillaPackageListByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackageList, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa package list successfully fetched!", villaPackageList)
}

func (h *VillaPackageListHandler) CreateVillaPackageList(c *gin.Context) {
	var req dto.VillaPackageListRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackageList, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Villa package list successfully created!", villaPackageList)
}

func (h *VillaPackageListHandler) UpdateVillaPackageList(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.VillaPackageListRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackageList, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa package list successfully updated!", villaPackageList)
}

func (h *VillaPackageListHandler) DeleteVillaPackageList(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Villa package list successfully deleted!", nil)
}
