// Package dto
package dto

type (
	RegisterRequest struct {
		Username string `json:"username" form:"username" binding:"required"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required,min=8"`
		Phone    string `json:"phone" form:"phone" binding:"required,min=11,max=11"`
	}

	LoginRequest struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	AuthResponse struct {
		User  UserResponse `json:"user"`
		Token string       `json:"token"`
	}
)
