// Package appError
package appError

import (
	"net/http"
)

type ServiceError struct {
	Message    string
	StatusCode int
	Errors     any
}

type errorsOpts struct {
	Errors any
}

type ErrorOption func(*errorsOpts)

func AddErrors(errors any) ErrorOption {
	return func(o *errorsOpts) { o.Errors = errors }
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) WithErrors(errs any) *ServiceError {
	e.Errors = errs
	return e
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

func Conflict(message string, opts ...ErrorOption) *ServiceError {
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
