package handlers

import (
	"net/http"

	"github.com/0-jagadeesh-0/chorvo/internal/api/v1/services"
	"github.com/0-jagadeesh-0/chorvo/internal/models/requests"
	"github.com/gin-gonic/gin"
)

// SignUpHandler handles user sign-up
// @Summary User Sign Up
// @Description Sign up a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body requests.CreateUserRequest true "User object"
// @Success 201 {object} responses.UserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /auth/signup [post]
func SignUpHandler(c *gin.Context) {
	var input requests.CreateUserRequest
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userResponse, err:= services.CreateUser(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

    c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": userResponse})
}

