package routes

import (
	"database/sql"

	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/controller"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/middleware"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/repository"
	"github.com/Gursevak56/food-delivery-platform/services/restaurant-service/services"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, db *sql.DB) {
	// repos
	restRepo := repository.NewRestaurantRepo(db)
	menuRepo := repository.NewMenuRepo(db)   // keep or implement separately
	orderRepo := repository.NewOrderRepo(db) // keep or implement separately

	// services
	restSvc := services.NewRestaurantService(restRepo)
	menuSvc := services.NewMenuService(menuRepo)
	orderSvc := services.NewOrderService(orderRepo, db)

	// controllers
	restC := controller.NewRestaurantController(restSvc)
	menuC := controller.NewMenuController(menuSvc)
	orderC := controller.NewOrderController(orderSvc)

	// restaurant routes
	rest := r.Group("/restaurants")
	{
		// public
		rest.GET("/", restC.GetAll)
		rest.GET("/:id", restC.Get)

		// protected - require auth
		auth := rest.Group("/")
		auth.Use(middleware.AuthRequired())

		auth.POST("/", restC.Create)
		auth.PUT("/:id", restC.Update)
		auth.DELETE("/:id", restC.Delete)
		auth.GET("/:id/qr/:table", restC.GenerateQR)

		// hours
		auth.POST("/:id/hours", restC.CreateHour)
		auth.GET("/:id/hours", restC.GetHours)
		auth.PUT("/:id/hours/:hours_id", restC.UpdateHour)
		auth.DELETE("/:id/hours/:hours_id", restC.DeleteHour)

		// tables (QR)
		auth.POST("/:id/tables", restC.CreateTable)
		auth.GET("/:id/tables", restC.ListTables)
		auth.PUT("/:id/tables/:table_id", restC.UpdateTable)
		auth.DELETE("/:id/tables/:table_id", restC.DeleteTable)
	}

	// keep menu & order endpoints wiring if implemented elsewhere
	rest.POST("/categories", menuC.CreateCategory)
	rest.GET("/:id/categories", menuC.GetCategories)
	rest.POST("/:id/menu/items", menuC.CreateMenuItem)
	rest.GET("/:id/menu/items", menuC.GetMenuItems)

	// orders / simple wiring example - implement order controller in order service file
	r.POST("/orders", orderC.PlaceOrder)
	r.GET("/orders/:id/status", orderC.GetStatus)
	r.PUT("/orders/:id/status", orderC.UpdateStatus)
}
