package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepositories interface {
	GetUserByEmail(db *gorm.DB, email string) (*model.UserModel, error)
	GetAllUser(db *gorm.DB) ([]*model.UserModel, error)
	CreateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error)
	UpdateUser(db *gorm.DB, user *model.UserModel) (*model.UserModel, error)
}

type UserServices interface {
	Login(req *request.ReqLogin) (*response.ResUser, error)
	Register(req *request.ReqCreateUser) (*response.ResUser, error)
}
type AuthHandlers interface {
	Logout(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}
