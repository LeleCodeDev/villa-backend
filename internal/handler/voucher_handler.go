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

type VoucherHandler struct {
	service *service.VoucherService
}

func NewVoucherHandler(service *service.VoucherService) *VoucherHandler {
	return &VoucherHandler{service: service}
}

func (h *VoucherHandler) GetAllVouchers(c *gin.Context) {
	var query dto.VoucherQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	vouchers, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {

		response.Success(c, http.StatusOK, "All vouchers successfully fetched!", vouchers)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All vouchers successfully fetched!", vouchers, pagination)
}

func (h *VoucherHandler) GetVoucherByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	voucher, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Voucher successfully fetched!", voucher)
}

func (h *VoucherHandler) ValidateVoucher(c *gin.Context) {
	var req dto.VoucherValidateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	currentUser := c.MustGet("user").(model.User)
	ctx := c.Request.Context()

	voucher, err := h.service.Validate(ctx, currentUser, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Voucher successfully validated!", voucher)
}

func (h *VoucherHandler) CreateVoucher(c *gin.Context) {
	var req dto.VoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	voucher, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Voucher successfully created!", voucher)
}

func (h *VoucherHandler) UpdateVoucher(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.VoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	voucher, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Voucher successfully updated!", voucher)
}

func (h *VoucherHandler) DeleteVoucher(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Voucher successfully updated!", nil)
}
