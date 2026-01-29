package userrepositories

import (
	"ecommerce-system/internal/models"

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

func (userRepo *UserRepositoryImpl) checkErrMysql(err error) error {
	if models.IsInternalErrMysql(err) {
		return err
	}
	if models.ForeignKeyErr(err) {
		return models.ErrRoleNotFound
	}
	if models.IsDuplicateKeyError(err) {
		return models.ErrDuplicateEmail
	}
	return models.ErrUserNotFound
}
func (userRepo *UserRepositoryImpl) GetUserPasswordByEmail(email string) (string, error) {
	var password string

	if err := userRepo.DB.Model(&models.UserModel{}).Select("password").Where("email=?", email).Scan(&password).Error; err != nil {
		return "", userRepo.checkErrMysql(err)
	}

	if password == "" {
		return "", models.ErrUserNotFound
	}
	return password, nil
}
func (userRepo *UserRepositoryImpl) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel

	if err := userRepo.DB.Preload("Role").Preload("Address").Where("email=?", email).Take(&user).Error; err != nil {
		return nil, userRepo.checkErrMysql(err)
	}

	return &user, nil
}
func (userRepo *UserRepositoryImpl) GetAllUser() ([]*models.UserModel, error) {
	var users []*models.UserModel

	if err := userRepo.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepositoryImpl) CreateUser(user *models.UserModel) (*models.UserModel, error) {

	if err := userRepo.DB.Create(user).Error; err != nil {
		return nil, userRepo.checkErrMysql(err)
	}

	return userRepo.GetUserByEmail(user.Email)
}

func (userRepo *UserRepositoryImpl) UpdateUser(user *models.UserModel) (*models.UserModel, error) {
	result := userRepo.DB.Model(&models.UserModel{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return nil, userRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected < 1 {
		return nil, models.ErrNoRowsAffected
	}

	return user, nil
}
