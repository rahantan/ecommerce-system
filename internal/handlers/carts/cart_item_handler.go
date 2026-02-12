package carthandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	cartitemsservices "ecommerce-system/internal/services/carts"

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
func (cartHandler *CartItemHandlerImpl) getUserData(ctx *fiber.Ctx) (*response.ResUser, error) {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return nil, exceptions.ErrCustomUnauthorized
	}
	return &user, nil
}

func (cartHandler *CartItemHandlerImpl) GetAllUserCartItem(ctx *fiber.Ctx) error {
	user, err := cartHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	result, err := cartHandler.CartItemServices.GetAllUserCartItem(user.ID)
	if err != nil {
		return err
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
	user, err := cartHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqCreateOrUpdateCartItem
	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	result, err := cartHandler.CartItemServices.CreateOrUpdateCartItem(&body, user.ID)
	if err != nil {
		return err
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
	user, err := cartHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqDeleteCartItem

	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	if err := cartHandler.CartItemServices.DeleteCartItemsByIDs(body.CartIDs, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success delete cart item",
	})
}
