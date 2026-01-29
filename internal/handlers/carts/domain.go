package carthandlers

import "github.com/gofiber/fiber/v2"

type CartItemHandlers interface {
	AddCartItem(ctx *fiber.Ctx) error
	GetAllUserCartItem(ctx *fiber.Ctx) error
	DeleteCartItemById(ctx *fiber.Ctx) error
	// GetCartItemById(ctx *fiber.Ctx) error //bisa diskip sih keknya
}
