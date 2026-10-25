package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToUserModel(req dto.UserRequest, hashedPassword string) *model.User {
	return &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     model.RoleUser,
		Password: hashedPassword,
	}
}

func ToUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
