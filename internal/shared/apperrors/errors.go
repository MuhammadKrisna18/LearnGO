package apperrors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code    int
	Message string
	Details interface{}
}

// Error implements the standard error interface
func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Details)
	}
	return e.Message
}

// Predefined error helpers

// NewBadRequest creates a 400 Bad Request error
func NewBadRequest(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: detail,
	}
}

// NewUnauthorized creates a 401 Unauthorized error
func NewUnauthorized(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Details: detail,
	}
}

// NewForbidden creates a 403 Forbidden error
func NewForbidden(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Details: detail,
	}
}

// NewNotFound creates a 404 Not Found error
func NewNotFound(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Details: detail,
	}
}

// NewInternal creates a 500 Internal Server Error
func NewInternal(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Details: detail,
	}
}
