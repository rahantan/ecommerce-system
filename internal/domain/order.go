package domain

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllOrder(db *gorm.DB) ([]*model.OrderModel, error)
	GetAllOrderByUserID(db *gorm.DB, userID int64) ([]*model.OrderModel, error)
	GetOrderDetailsByID(db *gorm.DB, userID, orderID int64) (*model.OrderModel, error)
	GetOrderByID(db *gorm.DB, orderID int64) (*model.OrderModel, error)
	CreateOrder(db *gorm.DB, order *model.OrderModel) (*model.OrderModel, error)
	DeleteOrder(db *gorm.DB, orderID int64, userID int64) error
	UpdateStatusOrder(db *gorm.DB, orderID int64, statusID int64) error
	UpdateAll(db *gorm.DB, order *model.OrderModel) error
}

type MidtransGateWay interface {
	CreateMidtrans(order *model.OrderModel) (*snap.Response, error)
}
type PaymentRepository interface {
	GetPaymentByOrderID(db *gorm.DB, orderID int64) (*model.PaymentOrderModel, error)
	UpdateStatusPayment(db *gorm.DB, orderID int64, status string) error
	SavePayment(db *gorm.DB, peyment *model.PaymentOrderModel) error
}

type OrderUseCase interface {
	GetOrderDetails(userID, orderID int64) (*response.ResOrder, error)
	GetAllOrderByUserID(userID int64) ([]*response.ResOrder, error)
	GetAllOrder() ([]*response.ResOrder, error)
	// UpdateStatusOrder(orderID, statusOrder int64) error
	ReceiveOrder(orderID, userID int64) error
	ShipOrder(orderID int64) error

	GetUserPaymentByOrderID(orderID, userID int64) (*response.ResPayment, error)

	UpdateStatusPayment(orderID int64, statusTransaction string) error
}

type OrderHandler interface {
	WebHookMidtransNotif(ctx *fiber.Ctx) error
	GetAllOrder(ctx *fiber.Ctx) error
	GetUserOrders(ctx *fiber.Ctx) error
	GetOrderDetails(ctx *fiber.Ctx) error
	ShipOrder(ctx *fiber.Ctx) error
	ReceiveOrder(ctx *fiber.Ctx) error
	GetUserPaymentByOrderID(ctx *fiber.Ctx) error
}
