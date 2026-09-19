package mapper

import (
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
)

func ToAuthResponse(user *model.User, token string) dto.AuthResponse {
	return dto.AuthResponse{
		User:  ToUserResponse(user),
		Token: token,
	}
}

func ToRegisteUserModel(req dto.RegisterRequest, hashedPassword string) *model.User {
	return &model.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     model.RoleUser,
	}
}
