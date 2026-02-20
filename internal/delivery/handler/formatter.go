package handler

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

func GetUserData(ctx *fiber.Ctx) (*response.ResUser, error) {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return nil, pkg.ErrCustomUnauthorized
	}

	return &user, nil
}
func SuccessResponse(ctx *fiber.Ctx, statusCode int, msg string, data any) error {
	return ctx.Status(statusCode).JSON(response.ResponseStandard{
		Success: true,
		Message: msg,
		Data:    data,
	})
}
