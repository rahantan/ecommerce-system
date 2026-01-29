package addresshandlers

import "github.com/gofiber/fiber/v2"

type AddressHandlers interface {
	CreateAddress(ctx *fiber.Ctx) error
	GetAllAddress(ctx *fiber.Ctx) error
	UpdateAddressByUserId(ctx *fiber.Ctx) error
	GetUserActiveAddress(ctx *fiber.Ctx) error
	// ActivateAddress(ctx *fiber.Ctx) error
}
