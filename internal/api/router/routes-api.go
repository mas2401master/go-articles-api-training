package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/internal/api/controller"
)

func SetudUsersRoutes(app *gin.Engine) {
	app.GET("/api/v1/users", controller.GetUsers)
	app.GET("/api//v1/users/:id", controller.GetUserById)
	app.POST("/api/v1/users", controller.CreateUser)
	app.PUT("/api/v1/users/:id", controller.UpdateUser)
	app.DELETE("/api/v1/users/:id", controller.DeleteUser)
}

func SetupItemRoutes(app *gin.Engine) {
	app.GET("/api/v1/items", controller.GetItems)
	app.GET("/api/v1/items/:id", controller.GetItemById)
	app.POST("/api/v1/items", controller.CreateItem)
	app.PUT("/api/v1/items/:id", controller.UpdateItem)
	app.DELETE("/api/v1/items/:id", controller.DeleteItem)
}

func SetupPromotionRoutes(app *gin.Engine) {
	//Promotion Routes
	app.GET("/api/v1/promotion", controller.GetPromotions)
	app.GET("/api/v1/promotion/:id", controller.GetPromotionById)
	app.POST("/api/v1/promotion", controller.CreatePromotion)
	app.PUT("/api/v1/promotion/:id", controller.UpdatePromotion)
	app.DELETE("/api/v1/promotion/:id", controller.DeletePromotion)
}

func SetudOrderRoutes(app *gin.Engine) {
	//Order  Routes
	app.GET("/api/v1/orders", controller.GetOrder)
	app.GET("/api/v1/orders/:id", controller.GetOrderById)
	app.POST("/api/v1/orders", controller.CreateOrder)
	app.PUT("/api/v1/orders/:id", controller.UpdateOrder)
	app.DELETE("/api/v1/orders/:id", controller.DeleteOrder)
}

func SetudOrderItemRoutes(app *gin.Engine) {
	//Order  Routes
	app.GET("/api/v1/orders/details/:id", controller.GetOrderItem)
	app.GET("/api/v1/orders/detail/:id", controller.GetOrderItemById)
	app.POST("/api/v1/orders/detail/:id", controller.CreateOrderItem)
	app.PUT("/api/v1/orders/detail/:id", controller.UpdateOrderItem)
	app.DELETE("/api/v1/orders/detail/:id", controller.DeleteOrderItem)
}
