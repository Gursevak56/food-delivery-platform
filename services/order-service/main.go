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
	InitDB()
	// Health check endpoint
	// r.GET("/health", func(c *gin.Context) {
	// 	healthCheckHandler(database)
	// 	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// })

	r.GET("/order/:id", getOrder)
	r.PUT("/order/:id", updateOrder)
	r.DELETE("/order/:id", deleteOrder)

	r.Run(":8081")
}

func getOrder(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "example"})
}

func updateOrder(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

func deleteOrder(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
