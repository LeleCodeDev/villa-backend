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

type VillaRuleHandler struct {
	service *service.VillaRuleService
}

func NewVillaRuleHandler(service *service.VillaRuleService) *VillaRuleHandler {
	return &VillaRuleHandler{service: service}
}

func (h *VillaRuleHandler) GetAllvillaRules(c *gin.Context) {
	var query dto.VillaRuleQuery

	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	villaRules, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {
		response.Success(c, http.StatusOK, "All villa rules successfully fetched!", villaRules)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All villa rules successfully fetched!", villaRules, pagination)
}

func (h *VillaRuleHandler) GetVillaRuleByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaRule, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa rule successfully fetched!", villaRule)
}

func (h *VillaRuleHandler) CreateVillaRule(c *gin.Context) {
	var req dto.VillaRuleRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaRule, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Villa rule successfully created!", villaRule)
}

func (h *VillaRuleHandler) UpdateVillaRule(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.VillaRuleRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	villaRule, err := h.service.Update(ctx, id, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Villa rule successfully updated!", villaRule)
}

func (h *VillaRuleHandler) DeleteVillaRule(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Villa rule successfully deleted!", nil)
}
