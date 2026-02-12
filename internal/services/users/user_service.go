package userservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	userrepositories "ecommerce-system/internal/repositories/users"
	"ecommerce-system/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	userrepositories.UserRepositories
	*gorm.DB
}

func NewUserService(user userrepositories.UserRepositories, db *gorm.DB) UserServices {
	return &UserServiceImpl{
		UserRepositories: user,
		DB:               db,
	}
}
func (userService *UserServiceImpl) loadUserRes(userMdl *models.UserModel) *response.ResUser {
	addresses := []response.ResAddress{}
	for _, address := range userMdl.Address {
		addresses = append(addresses, response.ResAddress{
			ID:      address.ID,
			City:    address.City,
			Address: address.Address,
		})
	}

	return &response.ResUser{
		ID:    userMdl.ID,
		Name:  userMdl.Name,
		Email: userMdl.Email,
		Phone: userMdl.Phone,
		Role: response.ResRole{
			ID:    userMdl.Role.ID,
			Title: userMdl.Role.Title,
		},
		Addresses: addresses,
		CreatedAt: userMdl.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: userMdl.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (userService *UserServiceImpl) Login(req *request.ReqLogin) (*response.ResUser, error) {
	user, err := userService.UserRepositories.GetUserByEmail(userService.DB, req.Email)
	if err != nil {
		return nil, exceptions.ErrCustomLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, exceptions.ErrCustomLogin
	}

	return userService.loadUserRes(user), nil
}
func (userService *UserServiceImpl) Register(request *request.ReqCreateUser) (*response.ResUser, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, exceptions.NewError(exceptions.KindInternal, err.Error(), nil)
	}

	request.Password = string(hash)
	result, err := userService.UserRepositories.CreateUser(userService.DB, &models.UserModel{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		RoleID:   request.RoleID,
	})
	if err != nil {
		return nil, utils.MappingError(err)
	}
	return userService.loadUserRes(result), nil
}
