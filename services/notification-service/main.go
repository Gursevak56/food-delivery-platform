package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment")
	}
}

func main() {
	r := gin.Default()
	database := InitDB()
	defer database.Close()
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/notifications", createNotification)
	r.GET("/notifications/:id", getNotification)
	r.PUT("/notifications/:id", updateNotification)
	r.DELETE("/notifications/:id", deleteNotification)

	r.Run(":8080")
}

func createNotification(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusCreated, gin.H{"message": "notification created"})
}

func getNotification(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "example"})
}

func updateNotification(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "notification updated"})
}

func deleteNotification(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
