package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthCheckHandler returns a Gin handler that pings the DB.
func healthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
