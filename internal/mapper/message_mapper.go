package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToMessageResponse(message *model.Message) dto.MessageResponse {
	return dto.MessageResponse{
		ID:        message.ID,
		Firstname: message.Firstname,
		Lastname:  message.Lastname,
		Email:     message.Email,
		Message:   message.Message,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func ToMessageModel(req dto.MessageRequest) *model.Message {
	return &model.Message{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Message:   req.Message,
	}
}

func UpdateMessageModel(message *model.Message, req dto.MessageRequest) {
	message.Firstname = req.Firstname
	message.Lastname = req.Lastname
	message.Email = req.Email
	message.Message = req.Message
}
