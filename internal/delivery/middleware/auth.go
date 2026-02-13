package middleware

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

func Authorization(allowedRole int64) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals("user").(response.ResUser)
		if !ok {
			return pkg.ErrCustomForbidden
		}

		if user.Role.ID != allowedRole {
			return pkg.ErrCustomForbidden
		}
		return ctx.Next()
	}
}
