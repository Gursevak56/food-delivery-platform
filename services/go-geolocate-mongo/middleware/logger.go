package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s - %s %s %d %s\n", clientIP, method, path, status, latency)
	}
}
