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
	if validationErr := ValidationError(err); validationErr != nil {
		return responseErrCustom(ctx, validationErr)
	}

	errCustom := err.(*ErrorCustom)
	if errCustom.StatusCode != http.StatusInternalServerError {
		return responseErrCustom(ctx, errCustom)
	}

	fmt.Println("error unexpected: ", errCustom.Errors)
	return ctx.Status(errCustom.StatusCode).JSON(response.ResponseStandard{
		Success: false,
		Message: errCustom.Message,
	})

}
