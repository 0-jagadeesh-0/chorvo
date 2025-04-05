package services

import (
	"github.com/0-jagadeesh-0/chorvo/internal/models/entities"
	"github.com/0-jagadeesh-0/chorvo/internal/models/requests"
	"github.com/0-jagadeesh-0/chorvo/internal/models/responses"
	repository "github.com/0-jagadeesh-0/chorvo/internal/repositories"
)

func CreateUser(userRequest *requests.CreateUserRequest) (*responses.UserResponse, error) {

	user := &entities.User{
        Name:     userRequest.Name,
        Email:    userRequest.Email,
        Password: userRequest.Password,
		Role: userRequest.Role,
    }

    createdUser, err := repository.CreateUser(user)
	if err != nil {
        return nil, err
    }

	userResponse := responses.UserResponse{
		ID:  createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  createdUser.Role,
	}

	return &userResponse, nil
}