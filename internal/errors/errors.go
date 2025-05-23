// Package errors - errors/errors.go
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ValidationError     ErrorType = "validation_error"
	NotFoundError       ErrorType = "not_found"
	ConflictError       ErrorType = "conflict"
	UnauthorizedError   ErrorType = "unauthorized"
	ForbiddenError      ErrorType = "forbidden"
	InternalError       ErrorType = "internal_error"
	BadRequestError     ErrorType = "bad_request"
	UnprocessableEntity ErrorType = "unprocessable_entity"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type APIError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(errorType ErrorType, message string, err error) *APIError {
	return &APIError{
		Type:       errorType,
		Message:    message,
		StatusCode: getStatusCode(errorType),
		Err:        err,
	}
}

func getStatusCode(errorType ErrorType) int {
	switch errorType {
	case ValidationError, BadRequestError:
		return http.StatusBadRequest
	case NotFoundError:
		return http.StatusNotFound
	case ConflictError:
		return http.StatusConflict
	case UnauthorizedError:
		return http.StatusUnauthorized
	case ForbiddenError:
		return http.StatusForbidden
	case InternalError:
		return http.StatusInternalServerError
	case UnprocessableEntity:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func NewValidationError(message string, err error) *APIError {
	return NewAPIError(ValidationError, message, err)
}

func NewNotFoundError(resource string, err error) *APIError {
	message := fmt.Sprintf("%s not found", resource)
	return NewAPIError(NotFoundError, message, err)
}

func NewConflictError(message string, err error) *APIError {
	return NewAPIError(ConflictError, message, err)
}

func NewUnauthorizedError(message string, err error) *APIError {
	if message == "" {
		message = "unauthorized access"
	}
	return NewAPIError(UnauthorizedError, message, err)
}

func NewForbiddenError(message string, err error) *APIError {
	if message == "" {
		message = "access forbidden"
	}
	return NewAPIError(ForbiddenError, message, err)
}

func NewInternalError(message string, err error) *APIError {
	if message == "" {
		message = "internal server error"
	}
	return NewAPIError(InternalError, message, err)
}

func NewBadRequestError(message string, err error) *APIError {
	return NewAPIError(BadRequestError, message, err)
}

func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

func GetAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}

func ToAPIError(err error) *APIError {
	if apiErr, ok := GetAPIError(err); ok {
		return apiErr
	}
	return NewInternalError("internal server error", err)
}
