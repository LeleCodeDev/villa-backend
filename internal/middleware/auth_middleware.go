package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/jwt"
	"github.com/lelecodedev/villa-backend/pkg/response"
)

func extractToken(c *gin.Context) *string {
	bearerToken := c.GetHeader("Authorization")
	if token, ok := strings.CutPrefix(bearerToken, "Bearer "); ok {
		return &token
	}

	return nil
}

func AuthMiddleware(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == nil {
			response.Error(c, http.StatusUnauthorized, "Authorization header is required!", nil)
			c.Abort()
			return
		}

		claims, err := jwt.ExtractToken(*tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token!", nil)
			c.Abort()
			return
		}

		userID, err := jwt.ExtractUserID(claims)
		if err != nil {
			response.HandleServiceError(c, err)
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Server errors", nil)
			c.Abort()
			return
		}

		c.Set("user", *user)
		c.Next()
	}
}
