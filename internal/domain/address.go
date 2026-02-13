package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AddressRepositories interface {
	GetAddressById(db *gorm.DB, addressID int64) (*model.AddressModel, error)
	GetUserAddressActive(db *gorm.DB, userId int64) (*model.AddressModel, error)
	GetAllAddress(db *gorm.DB, userId int64) ([]*model.AddressModel, error)
	UpdateAddressByUserId(db *gorm.DB, address *model.AddressModel) (*model.AddressModel, error)
	CreateAddress(db *gorm.DB, address *model.AddressModel) (*model.AddressModel, error)
	CheckNotFoundForUpdate(db *gorm.DB, addressID int64) error
	DeActivate(db *gorm.DB, userID int64) error
}

type AddressServices interface {
	GetUserAddressActive(userId int64) (*response.ResAddress, error)
	GetAllAddress(userid int64) ([]*response.ResAddress, error)
	CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error)
	UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error)
}

type AddressHandlers interface {
	CreateAddress(ctx *fiber.Ctx) error
	GetAllAddress(ctx *fiber.Ctx) error
	UpdateAddressByUserId(ctx *fiber.Ctx) error
	GetUserActiveAddress(ctx *fiber.Ctx) error
}
