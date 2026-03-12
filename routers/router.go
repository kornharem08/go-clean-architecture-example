package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/example/clean-architecture/internal/address"
	"github.com/example/clean-architecture/internal/middleware"
	"github.com/example/clean-architecture/internal/user"
)

// SetupRoutes initializes all routes for the application
func SetupRoutes(router *gin.Engine, userService user.IService, addressService address.IService) {
	// Apply global middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers
	// Inject multiple usecases into user handler via NewHandlerWithUsecases
	userHandler := user.NewHandlerWithUsecases(userService, addressService)
	addressHandler := address.NewHandler(addressService)

	// User routes
	usersGroup := router.Group("/api/users")
	{
		usersGroup.GET("", userHandler.ListUsers)
		usersGroup.POST("", userHandler.CreateUser)
		usersGroup.GET("/:id", userHandler.GetUser)
		usersGroup.PUT("/:id", userHandler.UpdateUser)
		usersGroup.DELETE("/:id", userHandler.DeleteUser)
	}

	// Address routes
	addressesGroup := router.Group("/api/addresses")
	{
		addressesGroup.GET("", addressHandler.ListAddresses)
		addressesGroup.POST("", addressHandler.CreateAddress)
		addressesGroup.GET("/:id", addressHandler.GetAddress)
		addressesGroup.PUT("/:id", addressHandler.UpdateAddress)
		addressesGroup.DELETE("/:id", addressHandler.DeleteAddress)
	}

	// User's addresses route
	userAddressesGroup := router.Group("/api/users/:user_id/addresses")
	{
		userAddressesGroup.GET("", addressHandler.GetAddressesByUser)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}
