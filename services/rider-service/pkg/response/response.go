package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, status int, message string, data any, meta any) {
	payload := gin.H{
		"success": true,
		"message": message,
		"data":    data,
	}
	if meta != nil {
		payload["meta"] = meta
	}
	c.JSON(status, payload)
}

func Error(c *gin.Context, status int, message, code string, errors any) {
	payload := gin.H{
		"success":    false,
		"message":    message,
		"error_code": code,
	}
	if errors != nil {
		payload["errors"] = errors
	}
	c.JSON(status, payload)
}
