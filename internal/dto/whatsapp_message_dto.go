package dto

type SendWhatsAppRequest struct {
	Phone   string `json:"phone" form:"phone" binding:"required,min=11,max=13"`
	Message string `json:"message" binding:"required"`
}
