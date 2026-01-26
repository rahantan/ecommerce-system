package categoryhandlers

import "github.com/gofiber/fiber/v2"

type CategoryHandlers interface {
	CreateCategory(ctx *fiber.Ctx) error
	GetAllCategory(ctx *fiber.Ctx) error
	UpdateCategoryById(ctx *fiber.Ctx) error
	GetCategoryById(ctx *fiber.Ctx) error
}
