package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderRepositories interface {
	CreateOrder(db *gorm.DB, order *model.OrderModel) error
	GetAllOrder(db *gorm.DB, userID int64) ([]*model.OrderModel, error)
	GetAllOrderItem(db *gorm.DB, orderID, userID int64) ([]*model.OrderItemModel, error)
	DeleteOrder(db *gorm.DB, orderID int64, userID int64) error
}

type CheckOutRepositories interface {
	CheckOut(db *gorm.DB, checkout *model.CheckoutModel) error
	CheckOutConfirm(db *gorm.DB, checkout *model.CheckoutModel, userID int64) error
	GetLastDraftCheckOut(db *gorm.DB, userID int64) (*model.CheckoutModel, error)
	UpdateStatusLastCheckOut(db *gorm.DB, status string, userID int64) error
}

type OrderServices interface {
	CheckOut(req *request.ReqCheckout, userID int64) error
	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) error
	GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error)
}

type OrderHandlers interface {
	CheckOut(ctx *fiber.Ctx) error
	CheckOutConfirm(ctx *fiber.Ctx) error
	GetLastDraftCheckOut(ctx *fiber.Ctx) error
}
