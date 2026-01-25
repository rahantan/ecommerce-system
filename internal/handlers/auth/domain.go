package authhandlers

import "github.com/gofiber/fiber/v2"

type AuthHandlers interface {
	Logout(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}
