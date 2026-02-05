package orderhandlers

import "github.com/gofiber/fiber/v2"

type OrderHandlers interface {
	CheckOut(ctx *fiber.Ctx) error
	CheckOutConfirm(ctx *fiber.Ctx) error
	GetLastDraftCheckOut(ctx *fiber.Ctx) error
}

// CheckOut(req request.ReqCheckout, userID int64) error
// 	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) error
