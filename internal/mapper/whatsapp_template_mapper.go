package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToWhatsappTemplateResponse(whatsappTemplate *model.WhatsappTemplate) dto.WhatsappTemplateResponse {
	return dto.WhatsappTemplateResponse{
		ID:        whatsappTemplate.ID,
		Text:      whatsappTemplate.Text,
		CreatedAt: whatsappTemplate.CreatedAt,
		UpdatedAt: whatsappTemplate.UpdatedAt,
	}
}

func ToWhatsappTemplateModel(req dto.WhatsappTemplateRequest) *model.WhatsappTemplate {
	return &model.WhatsappTemplate{
		Text: req.Text,
	}
}

func UpdateWhatsappTemplateModel(whatsappTemplate *model.WhatsappTemplate, req dto.WhatsappTemplateRequest) {
	whatsappTemplate.Text = req.Text
}
