package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/service"
	"github.com/lelecodedev/villa-backend/internal/util"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

type VillaPackageHandler struct {
	service *service.VillaPackageService
}

func NewVillaPackageHandler(service *service.VillaPackageService) *VillaPackageHandler {
	return &VillaPackageHandler{service: service}
}

func (h *VillaPackageHandler) GetAllVillaPackages(c *gin.Context) {
	ctx := c.Request.Context()

	villaPackages, err := h.service.GetAll(ctx)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "All villa package successfully fetched!", villaPackages)
}

func (h *VillaPackageHandler) GetVillaPackageByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackage, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa package successfully fetched!", villaPackage)
}

func (h *VillaPackageHandler) CreateVillaPackage(c *gin.Context) {
	var req dto.VillaPackageRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackage, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Villa package successfully created!", villaPackage)
}

func (h *VillaPackageHandler) UpdateVillaPackage(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.VillaPackageRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaPackage, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa package successfully updated!", villaPackage)
}

func (h *VillaPackageHandler) DeleteVillaPackage(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Villa package successfully deleted!", nil)
}
