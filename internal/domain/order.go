package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(db *gorm.DB, order *model.OrderModel) (*model.OrderModel, error)
	GetAllOrder(db *gorm.DB, userID int64) ([]*model.OrderModel, error)
	GetAllOrderItem(db *gorm.DB, orderID, userID int64) ([]*model.OrderItemModel, error)
	DeleteOrder(db *gorm.DB, orderID int64, userID int64) error
	UpdateStatusOrder(db *gorm.DB, orderID int64, statusID int64) error
	GetOrderByID(db *gorm.DB, orderID int64) (*model.OrderModel, error)
}

type CheckOutRepository interface {
	CheckOut(db *gorm.DB, checkout *model.CheckoutModel) error
	CheckOutConfirm(db *gorm.DB, checkout *model.CheckoutModel, userID int64) error
	GetLastDraftCheckOut(db *gorm.DB, userID int64) (*model.CheckoutModel, error)
	UpdateStatusLastCheckOut(db *gorm.DB, status string, userID int64) error
}

type MidtransGateWay interface {
	CreateMidtrans(order *model.OrderModel) (*snap.Response, error)
}

type OrderUseCase interface {
	CheckOut(req *request.ReqCheckout, userID int64) error
	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error)
	GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error)
	UpdateStatusOrder(orderID, statusOrder int64) error
}

type OrderHandler interface {
	CheckOut(ctx *fiber.Ctx) error
	CheckOutConfirm(ctx *fiber.Ctx) error
	GetLastDraftCheckOut(ctx *fiber.Ctx) error
	WebHookMidtransNotif(ctx *fiber.Ctx) error
}
