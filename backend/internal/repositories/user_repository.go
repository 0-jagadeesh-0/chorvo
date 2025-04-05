package repository

import (
	"github.com/0-jagadeesh-0/chorvo/database"
	"github.com/0-jagadeesh-0/chorvo/internal/models/entities"
)

func CreateUser(user *entities.User) (*entities.User, error) {
    err := database.DB.Create(user).Error
    if err != nil {
        return nil, err
    }
    return user, nil
}

func GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}
