// Package util
package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	appError "github.com/lelecodedev/villa-backend/internal/apperror"
)

func GetParamsID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, appError.BadRequest("Invalid ID")
	}

	if id < 0 {
		return 0, appError.BadRequest("Invalid ID")
	}

	return uint(id), nil
}
