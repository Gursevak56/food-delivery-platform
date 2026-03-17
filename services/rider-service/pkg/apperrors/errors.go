package apperrors

import "net/http"

type AppError struct {
	StatusCode int               `json:"-"`
	Code       string            `json:"error_code"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, code, message string) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message}
}

func Validation(errors map[string]string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_FAILED",
		Message:    "request validation failed",
		Errors:     errors,
	}
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(code, message string) *AppError {
	return New(http.StatusNotFound, code, message)
}

func Conflict(code, message string) *AppError {
	return New(http.StatusConflict, code, message)
}
