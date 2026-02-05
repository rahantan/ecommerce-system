package orderhandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	orderservices "ecommerce-system/internal/services/orders"
	"ecommerce-system/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandlerImpl struct {
	orderservices.OrderServices
	*validator.Validate
}

func NewOrderHandler(order orderservices.OrderServices, v *validator.Validate) OrderHandlers {
	return &OrderHandlerImpl{
		OrderServices: order,
		Validate:      v,
	}
}
func (orderHandler *OrderHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
}
func (orderHandler *OrderHandlerImpl) CheckOut(ctx *fiber.Ctx) error {

	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return orderHandler.withMessage(exceptions.ErrCustomUnauthorized, "checkout failed")
	}

	var body request.ReqCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return orderHandler.withMessage(exceptions.ErrCustomInvalidPayload, "checkout failed")
	}

	if err := orderHandler.Validate.Struct(&body); err != nil {
		return orderHandler.withMessage(exceptions.ValidationError(err), "checkout failed")
	}

	if err := orderHandler.OrderServices.CheckOut(&body, user.ID); err != nil {
		return orderHandler.withMessage(err, "checkout failed")
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success checkout",
	})
}
func (orderHandler *OrderHandlerImpl) CheckOutConfirm(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return orderHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to confirm checkout")
	}
	var body request.ReqConfirmCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return orderHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to confirm checkout")
	}
	if err := orderHandler.Validate.Struct(&body); err != nil {
		return orderHandler.withMessage(exceptions.ValidationError(err), "failed to confirm checkout")
	}

	if err := orderHandler.OrderServices.CheckOutConfirm(&body, user.ID); err != nil {
		return orderHandler.withMessage(err, "failed to confirm checkout")
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success confirm checkout",
	})
}
func (orderHandler *OrderHandlerImpl) GetLastDraftCheckOut(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return orderHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to get last checkout")
	}
	result, err := orderHandler.OrderServices.GetLastDraftCheckOut(user.ID)
	if err != nil {

		return orderHandler.withMessage(err, "failed to get last checkout")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get last draft checkout",
		Data: map[string]any{
			"checkout_draft": result,
		},
	})

}
