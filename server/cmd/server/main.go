package main

import (
	"fmt"
	"log"
	"os"

	"github.com/0-jagadeesh-0/chorvo/config"
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/handlers"
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/middleware"
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/routes"
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/services"
	"github.com/0-jagadeesh-0/chorvo/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
}

func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	routes.SetupAuthRoutes(router, authHandler)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// TODO: Add organization routes
		// TODO: Add team routes
		// TODO: Add project routes
		// TODO: Add task routes
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 