package user

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/pkg"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (addressHandler *UserHandlerImpl) CreateAddress(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqCreateAddress

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := addressHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := addressHandler.UserUseCase.CreateAddress(&body, user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusCreated, "success create address", map[string]any{
		"addresses": result,
	})

}
func (addressHandler *UserHandlerImpl) GetAllAddress(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	result, err := addressHandler.UserUseCase.GetAllAddress(user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get addresses", map[string]any{
		"addresses": result,
	})
}

func (addressHandler *UserHandlerImpl) UpdateAddressByUserId(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqUpdateAddress

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	addressId, err := strconv.Atoi(ctx.Params("addressId"))
	if err != nil {
		return pkg.ErrCustomInvalidAddressId
	}

	if err = addressHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := addressHandler.UserUseCase.UpdateAddressByUserId(&body, int64(addressId), user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success update address", map[string]any{
		"addresses": result,
	})
}
func (addressHandler *UserHandlerImpl) GetUserActiveAddress(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	result, err := addressHandler.UserUseCase.GetUserAddressActive(user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get address", map[string]any{
		"addresses": result,
	})
}
