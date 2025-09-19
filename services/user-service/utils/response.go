package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status     string      `json:"status"`          // "success" or "error"
	StatusCode int         `json:"statusCode"`      // numeric HTTP status code
	Message    string      `json:"message"`         // human message
	Data       interface{} `json:"data,omitempty"`  // response payload
	Error      interface{} `json:"error,omitempty"` // error details (optional)
}

// Success helper: sends standardized success response
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:     "success",
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

// Error helper: sends standardized error response
func SendError(c *gin.Context, statusCode int, message string, errDetail interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:     "error",
		StatusCode: statusCode,
		Message:    message,
		Error:      errDetail,
	})
}
