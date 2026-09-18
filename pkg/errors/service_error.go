package errors

import (
	"net/http"
)

type ServiceError struct {
	Message    string
	StatusCode int
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NotFound(message string) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func AlreadyExist(message string) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func Conflict(message string) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func Unauthorized(message string) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func BadRequest(message string) *ServiceError {
	return &ServiceError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
