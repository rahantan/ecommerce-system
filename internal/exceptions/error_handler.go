package exceptions

import (
	"ecommerce-system/internal/dto/response"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func responseErrCustom(ctx *fiber.Ctx, err *ErrorCustom) error {

	return ctx.Status(err.GetStatusCode()).JSON(response.ResponseStandard{
		Success: false,
		Message: string(err.Message),
		Errors:  err.Errors,
	})
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	if errCustom, ok := err.(*ErrorCustom); ok {
		return responseErrCustom(ctx, errCustom)
	}

	fmt.Println("error unexpected: ", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(response.ResponseStandard{
		Success: false,
		Message: "internal server error",
	})

}
