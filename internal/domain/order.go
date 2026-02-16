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
	GetAllOrder(db *gorm.DB, userID int64) ([]*model.OrderModel, error)
	GetOrderDetailsByID(db *gorm.DB, userID, orderID int64) (*model.OrderModel, error)
	GetOrderByID(db *gorm.DB, orderID int64) (*model.OrderModel, error)
	CreateOrder(db *gorm.DB, order *model.OrderModel) (*model.OrderModel, error)
	DeleteOrder(db *gorm.DB, orderID int64, userID int64) error
	UpdateStatusOrder(db *gorm.DB, orderID int64, statusID int64) error
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
type PaymentRepository interface {
	GetPaymentByOrderID(db *gorm.DB, orderID int64) (*model.PaymentOrderModel, error)
	UpdateStatusPayment(db *gorm.DB, orderID int64, status string) error
	SavePayment(db *gorm.DB, peyment *model.PaymentOrderModel) error
}

type OrderUseCase interface {
	GetOrderDetails(userID, orderID int64) (*response.ResOrder, error)
	GetAllOrder(userID int64) ([]response.ResOrder, error)
	UpdateStatusOrder(orderID, statusOrder int64) error

	CheckOut(req *request.ReqCheckout, userID int64) error
	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error)
	GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error)

	UpdateStatusPayment(orderID int64, statusTransaction string) error
}

type OrderHandler interface {
	CheckOut(ctx *fiber.Ctx) error
	CheckOutConfirm(ctx *fiber.Ctx) error
	GetLastDraftCheckOut(ctx *fiber.Ctx) error
	WebHookMidtransNotif(ctx *fiber.Ctx) error
	GetAllOrder(ctx *fiber.Ctx) error
	GetOrderDetails(ctx *fiber.Ctx) error
}
