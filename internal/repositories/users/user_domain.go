package userrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type UserRepositories interface {
	GetUserByEmail(db *gorm.DB, email string) (*models.UserModel, error)
	GetAllUser(db *gorm.DB) ([]*models.UserModel, error)
	CreateUser(db *gorm.DB, user *models.UserModel) (*models.UserModel, error)
	UpdateUser(db *gorm.DB, user *models.UserModel) (*models.UserModel, error)

	// GetUserPasswordByEmail(db *gorm.DB, email string) (string, error)
}
