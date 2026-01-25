package userservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type UserService interface {
	Create(request *request.ReqCreateUser) (*response.ResUser, error)
	GetUserByEmail(email string) (*response.ResUser, error)

	//method khusus
	// GetUserModelByEmail(email string) (*models.UserModel, error)
	GetUserPasswordByEmail(email string) (string, error)
}
