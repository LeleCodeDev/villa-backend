package dto

type SendWhatsAppRequest struct {
	Phone   string `json:"phone" form:"phone" binding:"required,min=12,max=12"`
	Message string `json:"message" binding:"required"`
}
