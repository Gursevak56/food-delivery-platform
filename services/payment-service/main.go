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

	// Example user endpoints
	r.POST("/", initiatePayment)
	r.GET("/:id", getPaymentStatus)
	r.PUT("/:id", updatePayment)
	r.DELETE("/:id", cancelPayment)

	r.Run(":8082")
}

func initiatePayment(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusCreated, gin.H{"message": "payment created"})
}

func getPaymentStatus(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "example"})
}

func updatePayment(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "payment updated"})
}

func cancelPayment(c *gin.Context) {
	// TODO: business logic
	c.JSON(http.StatusOK, gin.H{"message": "payment canceled"})
}
