package userservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	userrepositories "ecommerce-system/internal/repositories/users"
)

type UserServiceImpl struct {
	userrepositories.UserRepositories
}

func NewUserService(user userrepositories.UserRepositories) UserService {
	return &UserServiceImpl{
		UserRepositories: user,
	}
}
func (user *UserServiceImpl) handleError(err error) error {
	return exceptions.CheckError(err)
}
func (user *UserServiceImpl) Create(request *request.ReqCreateUser) (*response.ResUser, error) {

	result, err := user.UserRepositories.CreateUser(&models.UserModel{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		RoleID:   request.RoleID,
	})
	if errCheck := user.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return user.loadUserRes(result), nil
}

func (user *UserServiceImpl) GetUserByEmail(email string) (*response.ResUser, error) {
	result, err := user.UserRepositories.GetUserByEmail(email)
	if errCheck := user.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return user.loadUserRes(result), nil
}

func (user *UserServiceImpl) GetUserPasswordByEmail(email string) (string, error) {
	password, err := user.UserRepositories.GetUserPasswordByEmail(email)

	if errCheck := user.handleError(err); errCheck != nil {
		return "", errCheck
	}

	return password, nil
}
func (user *UserServiceImpl) loadUserRes(userMdl *models.UserModel) *response.ResUser {
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
