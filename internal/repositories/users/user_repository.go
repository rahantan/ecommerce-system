package userrepositories

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositories {
	return &UserRepositoryImpl{
		DB: db,
	}
}
func (userRepo *UserRepositoryImpl) GetUserPasswordByEmail(email string) (string, error) {
	var password string
	err := userRepo.DB.Model(&models.UserModel{}).Select("password").Where("email=?", email).Scan(&password).Error
	if err != nil {
		return "", err
	}

	if password == "" {
		return "", exceptions.ErrUserNotFound
	}
	return password, nil
}
func (userRepo *UserRepositoryImpl) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	err := userRepo.DB.Preload("Role").Preload("Address").Where("email=?", email).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
func (userRepo *UserRepositoryImpl) GetAllUser() ([]*models.UserModel, error) {
	var users []*models.UserModel
	err := userRepo.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepo *UserRepositoryImpl) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	result := userRepo.DB.Create(user)
	if result.Error != nil {
		if exceptions.IsDuplicateKeyError(result.Error) {
			return nil, exceptions.ErrDuplicateEmail
		}
		return nil, result.Error
	}

	return userRepo.GetUserByEmail(user.Email)
}

func (userRepo *UserRepositoryImpl) UpdateUser(user *models.UserModel) (*models.UserModel, error) {
	result := userRepo.DB.Model(&models.UserModel{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, exceptions.ErrNoRowsAffected
	}
	return user, nil
}
