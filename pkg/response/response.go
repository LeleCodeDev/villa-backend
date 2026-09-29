package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lelecodedev/villa-backend/pkg/errors"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Errors  any    `json:"errors,omitempty"`
}

type PaginatedResponse[T any] struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       T              `json:"data"`
	Pagination PaginationData `json:"pagination"`
	Errors     any            `json:"errors,omitempty"`
}

type PaginationData struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func Success[T any](c *gin.Context, code int, message string, data T) {
	c.JSON(code, Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Paginated[T any](c *gin.Context, code int, message string, data T, pagination PaginationData) {
	c.JSON(code, PaginatedResponse[T]{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Error(c *gin.Context, code int, message string, err any) {
	c.JSON(code, Response[any]{
		Success: false,
		Message: message,
		Errors:  err,
	})
}

func HandleServiceError(c *gin.Context, err error) {
	if serviceError, ok := err.(*errors.ServiceError); ok {
		Error(c, serviceError.StatusCode, serviceError.Message, nil)
		return
	}

	Error(c, http.StatusInternalServerError, "Server error", err.Error())
}

func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		Error(c, http.StatusBadRequest, "Validation failed!", errors.GetValidationError(validationErrors))
		return
	}

	Error(c, http.StatusInternalServerError, "Server error", err.Error())
}
