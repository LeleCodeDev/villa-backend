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

type WhatsappTemplateHandler struct {
	service *service.WhatsappTemplateService
}

func NewWhatsappTemplateHandler(service *service.WhatsappTemplateService) *WhatsappTemplateHandler {
	return &WhatsappTemplateHandler{service: service}
}

func (h *WhatsappTemplateHandler) GetAllWhatsappTemplates(c *gin.Context) {
	var query dto.WhatsappTemplateQuery
	if err := c.ShouldBind(&query); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	query.SetDefault()
	ctx := c.Request.Context()

	whatsappTemplates, total, err := h.service.GetAll(ctx, query)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	if query.Unpage {

		response.Success(c, http.StatusOK, "All whatsapp templates successfully fetched!", whatsappTemplates)
		return
	}

	pagination := pagination.BuildPagination(query.Page, query.Size, total)
	response.Paginated(c, http.StatusOK, "All whatsapp templates successfully fetched!", whatsappTemplates, pagination)
}

func (h *WhatsappTemplateHandler) GetWhatsappTemplateByID(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	ctx := c.Request.Context()

	whatsappTemplate, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Whatsapp template successfully fetched!", whatsappTemplate)
}

func (h *WhatsappTemplateHandler) CreateWhatsappTemplate(c *gin.Context) {
	var req dto.WhatsappTemplateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	whatsappTemplate, err := h.service.Create(ctx, req)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Whatsapp template successfully created!", whatsappTemplate)
}

func (h *WhatsappTemplateHandler) UpdateWhatsappTemplate(c *gin.Context) {
	id, err := util.GetParamsID(c)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	var req dto.WhatsappTemplateRequest
	if err := c.ShouldBind(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	ctx := c.Request.Context()

	whatsappTemplate, err := h.service.Update(ctx, req, id)
	if err != nil {
		response.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Whatsapp template successfully updated!", whatsappTemplate)
}

func (h *WhatsappTemplateHandler) DeleteWhatsappTemplate(c *gin.Context) {
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

	response.Success[any](c, http.StatusOK, "Whatsapp template successfully updated!", nil)
}
