package routes

import (
	"github.com/gin-gonic/gin"

	"minimart/handlers"
	"minimart/middleware"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Public routes
	api.POST("/auth/signup", handlers.Signup)
	api.POST("/auth/login", handlers.Login)

	api.GET("/products", handlers.GetProducts)
	api.GET("/products/:id", handlers.GetProduct)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/auth/logout", handlers.Logout)
	protected.GET("/profile", handlers.GetProfile)

	// Cart routes
	protected.GET("/cart", handlers.GetCart)
	protected.POST("/cart", handlers.AddToCart)
	protected.PATCH("/cart/:id", handlers.UpdateCartItem)
	protected.DELETE("/cart/:id", handlers.RemoveCartItem)

	// Order routes
	protected.POST("/orders", handlers.PlaceOrder)
	protected.GET("/orders", handlers.GetOrders)
}
