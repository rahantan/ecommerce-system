package repository

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() domain.UserRepository {
	return &UserRepositoryImpl{}
}

func (userRepo *UserRepositoryImpl) checkErrMysql(err error) error {
	if model.IsInternalErrMysql(err) {
		return err
	}
	if model.ForeignKeyErr(err) {
		return model.ErrRoleNotFound
	}
	if model.IsDuplicateKeyError(err) {
		return model.ErrDuplicateEmail
	}
	return model.ErrUserNotFound
}

func (userRepo *UserRepositoryImpl) GetUserByEmail(db *gorm.DB, email string) (*model.UserModel, error) {
	var user model.UserModel

	if err := db.Preload("Role").Preload("Address").Where("email=?", email).Take(&user).Error; err != nil {
		return nil, userRepo.checkErrMysql(err)
	}

	return &user, nil
}
func (userRepo *UserRepositoryImpl) GetAllUser(db *gorm.DB) ([]*model.UserModel, error) {
	var users []*model.UserModel

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (userRepo *UserRepositoryImpl) CreateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error) {

	if err := db.Create(user).Error; err != nil {
		return nil, userRepo.checkErrMysql(err)
	}

	return userRepo.GetUserByEmail(db, user.Email)
}

func (userRepo *UserRepositoryImpl) UpdateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error) {
	result := db.Model(&model.UserModel{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return nil, userRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected < 1 {
		return nil, model.ErrNoRowsAffected
	}

	return user, nil
}
