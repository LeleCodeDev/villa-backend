package dto

import (
	"time"

	"github.com/lelecodedev/villa-backend/internal/model"
)

type (
	UserQuery struct {
		PaginationQuery
		SortQuery
		Unpage bool   `form:"unpage"`
		Search string `form:"search"`
	}

	UserRequest struct {
		Username string `json:"username" form:"username" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Phone    string `json:"phone" form:"phone" binding:"required,min=12,max=12"`
		Password string `json:"password" form:"password" binding:"required,min=8"`
	}

	UserResponse struct {
		ID        uint           `json:"id"`
		Username  string         `json:"username"`
		Email     string         `json:"email"`
		Phone     string         `json:"phone"`
		Role      model.UserRole `json:"role"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
	}
)

func (uq *UserQuery) SetDefault() {
	uq.setDefaultPagination()
	uq.SetDefaultSort(SortAsc)
}
