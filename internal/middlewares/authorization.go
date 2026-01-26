package middlewares

import (
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"

	"github.com/gofiber/fiber/v2"
)

func Authorization(allowedRole int64) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(response.ResUser)
		if !ok {
			return exceptions.ErrCustomForbidden
		}

		if user.Role.ID != allowedRole {
			return exceptions.ErrCustomForbidden
		}
		return ctx.Next()
	}
}
