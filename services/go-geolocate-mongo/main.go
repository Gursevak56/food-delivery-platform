package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go-geolocate-mongo/handlers"
	"go-geolocate-mongo/middleware"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment")
	}
}

func main() {
	InitMongo()
	InitDB()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SimpleCORS())

	// routes
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/detect/ip", handlers.DetectIP)
	r.POST("/reverse", handlers.ReverseGeocode)
	r.POST("/save", handlers.SaveLocation)
	r.GET("/locations", handlers.ListLocations)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s\n", port)
	r.Run(":" + port)
}
