package cart

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CartItemHandlerImpl struct {
	domain.CartUseCase
	*validator.Validate
}

func NewCartItemHandler(cartService domain.CartUseCase, v *validator.Validate) domain.CartHandler {
	return &CartItemHandlerImpl{
		CartUseCase: cartService,
		Validate:    v,
	}
}

func (cartHandler *CartItemHandlerImpl) GetAllUserCartItem(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	result, err := cartHandler.CartUseCase.GetAllUserCartItem(user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get cart items", map[string]any{
		"items": result,
	})

}

func (cartHandler *CartItemHandlerImpl) AddCartItem(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqAddCart
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := cartHandler.CartUseCase.AddCartItem(&body, user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success add cart item", map[string]any{
		"items": result,
	})

}

func (cartHandler *CartItemHandlerImpl) UpdateCartItemByID(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	cartID, err := ctx.ParamsInt("cartID")
	if err != nil {
		return pkg.ErrCustomInvalidCartId
	}

	var body request.ReqUpdateCartQty
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := cartHandler.CartUseCase.UpdateCartItemByID(&body, int64(cartID), user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success update cart item", map[string]any{
		"items": result,
	})

}

func (cartHandler *CartItemHandlerImpl) DeleteCartItemsByIDs(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqDeleteCartItem

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := cartHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	if err := cartHandler.CartUseCase.DeleteCartItemsByIDs(body.CartIDs, user.ID); err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success delete cart item", nil)

}
