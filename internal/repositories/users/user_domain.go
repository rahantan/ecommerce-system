package userrepositories

import (
	"ecommerce-system/internal/models"
)

type UserRepositories interface {
	GetUserByEmail(email string) (*models.UserModel, error)
	GetAllUser() ([]*models.UserModel, error)
	CreateUser(user *models.UserModel) (*models.UserModel, error)
	UpdateUser(user *models.UserModel) (*models.UserModel, error)

	GetUserPasswordByEmail(email string) (string, error)
}
