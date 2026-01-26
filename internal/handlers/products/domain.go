package producthandlers

import "github.com/gofiber/fiber/v2"

type ProductHandlers interface {
	CreateProduct(ctx *fiber.Ctx) error
	GetAllProduct(ctx *fiber.Ctx) error
	UpdateProductById(ctx *fiber.Ctx) error
	GetProductById(ctx *fiber.Ctx) error
}
