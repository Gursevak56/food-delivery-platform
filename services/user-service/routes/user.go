package routes

import (
	"github.com/Gursevak56/food-delivery-platform/services/user-service/controller"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	Controller *controller.UserController
}

func (route *UserRoute) RegisterRoutes(r *gin.Engine) {
	r.POST("/auth/send-otp", route.Controller.SendOTP)
	r.POST("/auth/verify-otp", route.Controller.VerifyOTP)

	users := r.Group("/users")
	{
		users.POST("/", route.Controller.CreateUser)
		users.POST("/login", route.Controller.LoginUser)
		users.GET("/", route.Controller.GetAllUsers)

		// user-specific routes (use :id consistently)
		users.GET("/:id", route.Controller.GetUserByID)
		users.PUT("/:id", route.Controller.UpdateUser)
		users.DELETE("/:id", route.Controller.DeleteUser)

		// protected addresses under users/:id/addresses
		// require auth first
		usersAuth := users.Group("/:id")
		usersAuth.Use(middleware.AuthRequired())
		{
			// create address (owner only)
			usersAuth.POST("/addresses", middleware.OwnerOnly(), route.Controller.CreateAddress)
			// get addresses for the user (owner or admin)
			usersAuth.GET("/addresses", middleware.OwnerOrAdmin(), route.Controller.GetAddresses)
		}
	}

	// addresses as independent resource (single address operations)
	addresses := r.Group("/addresses")
	addresses.Use(middleware.AuthRequired())
	{
		// get single address (auth)
		addresses.GET("/:id", route.Controller.GetAddressByID)
		// update address (owner only)
		addresses.PUT("/:id", middleware.OwnerOnly(), route.Controller.UpdateAddress)
		// delete (owner only)
		addresses.DELETE("/:id", middleware.OwnerOnly(), route.Controller.DeleteAddress)
	}
}
