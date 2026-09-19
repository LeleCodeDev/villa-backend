package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	UserResponse struct {
		ID        uint           `json:"id"`
		Username  string         `json:"username"`
		Email     string         `json:"email"`
		Phone     string         `json:"phone"`
		Role      model.UserRole `json:"role"`
		CreatedAt time.Time      `json:"createdAt"`
		UpdatedAt time.Time      `json:"updatedAt"`
	}
)
