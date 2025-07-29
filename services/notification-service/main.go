package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	database := InitDB()
	defer database.Close()
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Example user endpoints
	r.GET("/health", healthCheckHandler(database))
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}

func createUser(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func getUser(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "example"})
}

func updateUser(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func deleteUser(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
