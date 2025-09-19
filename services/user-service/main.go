package main

import (
	"context"
	"log"
	"os"

	"github.com/Gursevak56/food-delivery-platform/services/user-service/controller"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/repository"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/routes"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env silently (no error if file missing)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment")
	}
}

func main() {
	ctx := context.Background()
	// 1) Initialize the DB
	database := InitDB()
	defer database.Close()
	mongoClient, err := InitMongo(ctx)
	defer mongoClient.Disconnect(context.Background())

	if err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}
	otpCol := mongoClient.Database("food_app").Collection("otps")
	// ensure TTL/indexes (optional error handling)
	if err := repository.EnsureOTPIndexes(ctx, otpCol); err != nil {
		log.Printf("warning: EnsureOTPIndexes failed: %v", err)
	}

	// 3) Construct repositories & services, injecting OTPRepo into UserService
	userRepo := &repository.UserRepository{DB: database}
	otpRepo := repository.NewOTPRepo(otpCol) // <- OTP repo constructed here

	userSvc := &service.UserService{
		Repo:    userRepo,
		OTPRepo: otpRepo, // <- inject here
	}

	userCtrl := &controller.UserController{
		Service: userSvc,
	}

	// 4) Create router and register routes
	r := gin.Default()
	// r.GET("/health", healthCheckHandler(database))

	// wire route with controller
	UserRoute := &routes.UserRoute{
		Controller: userCtrl,
	}
	UserRoute.RegisterRoutes(r)

	// 5) Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}
