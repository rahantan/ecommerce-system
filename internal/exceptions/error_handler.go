package exceptions

import (
	"ecommerce-system/internal/dto/response"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func responseErrCustom(ctx *fiber.Ctx, err *ErrorCustom) error {

	return ctx.Status(err.StatusCode).JSON(response.ResponseStandard{
		Success: false,
		Message: err.Message,
		Errors:  err.Errors,
	})
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {

	errCustom, ok := err.(*ErrorCustom)

	if ok && errCustom.StatusCode != http.StatusInternalServerError {
		res := responseErrCustom(ctx, errCustom)
		return res
	}

	fmt.Println("error unexpected: ", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(response.ResponseStandard{
		Success: false,
		Message: "internal server error",
	})

}
