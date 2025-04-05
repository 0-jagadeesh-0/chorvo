package router

import (
	_ "github.com/0-jagadeesh-0/chorvo/docs"
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Set up routes
	v1 := router.Group("/api/v1")
	routes.RegisterAuthRoutes(v1)

	return router
}