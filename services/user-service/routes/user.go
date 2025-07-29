package routes

import (
	"github.com/Gursevak56/food-delivery-platform/services/user-service/controller"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	Controller *controller.UserController
}

func (route *UserRoute) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/users")
	users.POST("/", route.Controller.CreateUser)
	users.POST("/login", route.Controller.LoginUser)
	users.GET("/", route.Controller.GetAllUsers)
	users.GET("/:id", route.Controller.GetUserByID)
	users.PUT("/:id", route.Controller.UpdateUser)
	users.DELETE("/:id", route.Controller.DeleteUser)
}
