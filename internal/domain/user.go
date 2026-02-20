package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(db *gorm.DB, email string) (*model.UserModel, error)
	GetAllUser(db *gorm.DB) ([]*model.UserModel, error)
	CreateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error)
	UpdateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error)
}

type UserUseCase interface {
	Login(req *request.ReqLogin) (*response.ResUser, error)
	Register(req *request.ReqCreateUser) (*response.ResUser, error)

	GetUserAddressActive(userId int64) (*response.ResAddress, error)
	GetAllAddress(userid int64) ([]*response.ResAddress, error)
	CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error)
	UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error)
}
type UserHandler interface {
	Logout(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error

	CreateAddress(ctx *fiber.Ctx) error
	GetAllAddress(ctx *fiber.Ctx) error
	UpdateAddressByUserId(ctx *fiber.Ctx) error
	GetUserActiveAddress(ctx *fiber.Ctx) error
}
