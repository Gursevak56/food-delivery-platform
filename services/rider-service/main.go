package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	database := InitDB()
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		healthCheckHandler(database)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Example user endpoints
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8084")
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
