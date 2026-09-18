package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lelecodedev/villa-backend/pkg/errors"
)

func GetParamsID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errors.BadRequest("Invalid ID")
	}

	if id < 0 {
		return 0, errors.BadRequest("Invalid ID")
	}

	return uint(id), nil
}
