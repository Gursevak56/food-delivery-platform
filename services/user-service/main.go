package main

import (
	"log"

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
	// 1) Initialize the DB
	database := InitDB()
	defer database.Close()

	// 2) Create your router
	r := gin.Default()

	// 3) Use your healthCheckHandler
	r.GET("/health", healthCheckHandler(database))

	type UserMain struct {
		routes *routes.UserRoute
	}
	UserRoute := &routes.UserRoute{
		Controller: &controller.UserController{
			Service: &service.UserService{
				Repo: &repository.UserRepository{
					DB: database,
				},
			},
		},
	}
	UserRoute.RegisterRoutes(r)

	// 5) Run the server
	r.Run("0.0.0.0:8085")
}
