package usecase

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
	*gorm.DB
}

func NewUserUseCase(user domain.UserRepository, db *gorm.DB) domain.UserUseCase {
	return &UserUseCaseImpl{
		UserRepository: user,
		DB:             db,
	}
}
func (userUC *UserUseCaseImpl) loadUserRes(userMdl *model.UserModel) *response.ResUser {
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
