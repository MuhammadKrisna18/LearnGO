package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Details interface{}
}

func (e *AppError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Details)
	}
	return e.Message
}

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

func NewConflict(message string, details ...interface{}) *AppError {
	var detail interface{}
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Details: detail,
	}
}
