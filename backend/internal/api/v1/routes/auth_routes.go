package routes

import (
	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {

	
	authGroup := router.Group("/auth")

	authGroup.POST("/signup", handlers.SignUpHandler)

}