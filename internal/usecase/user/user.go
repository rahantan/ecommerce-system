package user

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCaseImpl struct {
	domain.UserRepository
	domain.AddressRepository
	*gorm.DB
}

func NewUserUseCase(user domain.UserRepository, address domain.AddressRepository, db *gorm.DB) domain.UserUseCase {
	return &UserUseCaseImpl{
		UserRepository:    user,
		DB:                db,
		AddressRepository: address,
	}
}

func (userUC *UserUseCaseImpl) Login(req *request.ReqLogin) (*response.ResUser, error) {
	user, err := userUC.UserRepository.GetUserByEmail(userUC.DB, req.Email)
	if err != nil {
		return nil, pkg.ErrCustomLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, pkg.ErrCustomLogin
	}

	return userUC.loadUserRes(user), nil
}

func (userUC *UserUseCaseImpl) Register(request *request.ReqCreateUser) (*response.ResUser, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, pkg.NewError(pkg.KindInternal, err.Error(), nil)
	}

	request.Password = string(hash)
	result, err := userUC.UserRepository.CreateUser(userUC.DB, &model.UserModel{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		RoleID:   request.RoleID,
	})
	if err != nil {
		return nil, pkg.MappingError(err)
	}
	return userUC.loadUserRes(result), nil
}
