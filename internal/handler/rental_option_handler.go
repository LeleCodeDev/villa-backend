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

type RentalOptionHandler struct {
	service *service.RentalOptionService
}

func NewRentalOptionHandler(service *service.RentalOptionService) *RentalOptionHandler {
	return &RentalOptionHandler{service: service}
}

func (h *RentalOptionHandler) GetAllRentalOptions(c *gin.Context) {
	var query dto.RentalOptionQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	rentalOptions, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {

		response.Success(c, http.StatusOK, "All reantal options successfully fetched!", rentalOptions)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All rental options successfully fetched!", rentalOptions, pagination)
}

func (h *RentalOptionHandler) GetRentalOptionByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	rentalOption, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Rental option successfully fetched!", rentalOption)
}

func (h *RentalOptionHandler) CreateRentalOption(c *gin.Context) {
	var req dto.RentalOptionRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	rentalOption, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Rental option successfully created!", rentalOption)
}

func (h *RentalOptionHandler) UpdateRentalOption(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.RentalOptionRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	rentalOption, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Rental option successfully updated!", rentalOption)
}

func (h *RentalOptionHandler) DeleteRentalOption(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Rental option successfully deleted!", nil)
}
