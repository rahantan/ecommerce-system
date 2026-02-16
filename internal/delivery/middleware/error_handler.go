package middleware

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/pkg"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func responseErrCustom(ctx *fiber.Ctx, err *pkg.ErrorCustom) error {

	return ctx.Status(err.GetStatusCode()).JSON(response.ResponseStandard{
		Success: false,
		Message: err.Message,
		Errors:  err.Errors,
	})
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	if errCustom, ok := err.(*pkg.ErrorCustom); ok {
		return responseErrCustom(ctx, errCustom)
	}

	fmt.Println("error unexpected: ", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(response.ResponseStandard{
		Success: false,
		Message: "internal server error",
	})

}
