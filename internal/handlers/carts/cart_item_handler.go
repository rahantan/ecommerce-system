package carthandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	cartitemsservices "ecommerce-system/internal/services/carts"
	"ecommerce-system/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CartItemHandlerImpl struct {
	cartitemsservices.CartItemServices
	*validator.Validate
}

func NewCartItemHandler(cartService cartitemsservices.CartItemServices, v *validator.Validate) CartItemHandlers {
	return &CartItemHandlerImpl{
		CartItemServices: cartService,
		Validate:         v,
	}
}
func (cartHandler *CartItemHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
}
func (cartHandler *CartItemHandlerImpl) GetAllUserCartItem(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return cartHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to get cart items")
	}

	result, err := cartHandler.CartItemServices.GetAllUserCartItem(user.ID)
	if err != nil {
		return cartHandler.withMessage(err, "failed to get cart items")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get cart items",
		Data: map[string]any{
			"items": result,
		},
	})
}
func (cartHandler *CartItemHandlerImpl) AddCartItem(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return cartHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to add cart item")
	}

	var body request.ReqCreateOrUpdateCartItem
	if err := ctx.BodyParser(&body); err != nil {
		return cartHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to add cart item")
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return cartHandler.withMessage(exceptions.ValidationError(err), "failed to add cart item")
	}

	result, err := cartHandler.CartItemServices.CreateOrUpdateCartItem(&body, user.ID)
	if err != nil {
		return cartHandler.withMessage(err, "failed to add cart item")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success add cart item",
		Data: map[string]any{
			"item": result,
		},
	})
}

func (cartHandler *CartItemHandlerImpl) DeleteCartItemsByIDs(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return cartHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to delete cart item")
	}
	var body request.ReqDeleteCartItem
	if err := ctx.BodyParser(&body); err != nil {
		return cartHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to delete cart item")
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return cartHandler.withMessage(exceptions.ValidationError(err), "failed to delete cart item")
	}

	if err := cartHandler.CartItemServices.DeleteCartItemsByIDs(body.CartIDs, user.ID); err != nil {
		return cartHandler.withMessage(err, "failed to delete cart item")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success delete cart item",
	})
}
