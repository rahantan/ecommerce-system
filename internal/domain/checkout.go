package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CheckOutRepository interface {
	CheckOut(db *gorm.DB, checkout *model.CheckoutModel) error
	CheckOutConfirm(db *gorm.DB, checkout *model.CheckoutModel, userID int64) error
	GetLastDraftCheckOut(db *gorm.DB, userID int64) (*model.CheckoutModel, error)
	UpdateStatusLastCheckOut(db *gorm.DB, status string, userID int64) error
}

type CheckOutUseCase interface {
	CheckOut(req *request.ReqCheckout, userID int64) error
	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error)
	GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error)
}

type CheckOutHandler interface {
	CheckOut(ctx *fiber.Ctx) error
	CheckOutConfirm(ctx *fiber.Ctx) error
	GetLastDraftCheckOut(ctx *fiber.Ctx) error
}
