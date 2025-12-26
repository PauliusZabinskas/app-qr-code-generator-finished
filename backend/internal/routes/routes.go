package routes

import (
	"gin-quickstart/internal/config"
	"gin-quickstart/internal/handlers"
	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/repositories"
	"gin-quickstart/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	wifiRepo := repositories.NewWifiRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	qrCodeService := services.NewQRCodeService()
	wifiService := services.NewWifiService(wifiRepo, qrCodeService, cfg.EncryptionKey)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	wifiHandler := handlers.NewWifiHandler(wifiService)
	adminHandler := handlers.NewAdminHandler(userRepo, wifiRepo)

	// API route group
	api := router.Group("/api")
	{
		// Public auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected WiFi routes
		wifi := api.Group("/wifi")
		wifi.Use(middleware.AuthMiddleware(authService))
		{
			wifi.GET("", wifiHandler.GetAll)
			wifi.POST("", wifiHandler.Create)
			wifi.GET("/:id", wifiHandler.GetByID)
			wifi.DELETE("/:id", wifiHandler.Delete)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/credentials", adminHandler.GetAllCredentials)
			admin.GET("/stats", adminHandler.GetStats)
		}
	}
}
