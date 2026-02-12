package orderhandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	orderservices "ecommerce-system/internal/services/orders"

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
func (orderHandler *OrderHandlerImpl) getUserData(ctx *fiber.Ctx) (*response.ResUser, error) {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return nil, exceptions.ErrCustomUnauthorized
	}

	return &user, nil
}
func (orderHandler *OrderHandlerImpl) CheckOut(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := orderHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	if err := orderHandler.OrderServices.CheckOut(&body, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success checkout",
	})
}
func (orderHandler *OrderHandlerImpl) CheckOutConfirm(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqConfirmCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := orderHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	if err := orderHandler.OrderServices.CheckOutConfirm(&body, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success confirm checkout",
	})
}
func (orderHandler *OrderHandlerImpl) GetLastDraftCheckOut(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	result, err := orderHandler.OrderServices.GetLastDraftCheckOut(user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get last draft checkout",
		Data: map[string]any{
			"checkout_draft": result,
		},
	})

}
