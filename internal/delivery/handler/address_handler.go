package handler

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"

	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AddressHandlerImpl struct {
	domain.AddressServices
	*validator.Validate
}

func NewAddressHandler(address domain.AddressServices, v *validator.Validate) domain.AddressHandlers {
	return &AddressHandlerImpl{
		AddressServices: address,
		Validate:        v,
	}
}
func (addressHandler *AddressHandlerImpl) getUserData(ctx *fiber.Ctx) (*response.ResUser, error) {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return nil, pkg.ErrCustomUnauthorized
	}
	return &user, nil
}

func (addressHandler *AddressHandlerImpl) CreateAddress(ctx *fiber.Ctx) error {

	user, err := addressHandler.getUserData(ctx)
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

	result, err := addressHandler.AddressServices.CreateAddress(&body, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success create address",
		Data: map[string]any{
			"address": result,
		},
	})
}
func (addressHandler *AddressHandlerImpl) GetAllAddress(ctx *fiber.Ctx) error {

	user, err := addressHandler.getUserData(ctx)
	if err != nil {
		return err
	}
	result, err := addressHandler.AddressServices.GetAllAddress(user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get addresses",
		Data: map[string]any{
			"addresses": result,
		},
	})
}

func (addressHandler *AddressHandlerImpl) UpdateAddressByUserId(ctx *fiber.Ctx) error {

	user, err := addressHandler.getUserData(ctx)
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

	result, err := addressHandler.AddressServices.UpdateAddressByUserId(&body, int64(addressId), user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success update address",
		Data: map[string]any{
			"address": result,
		},
	})
}
func (addressHandler *AddressHandlerImpl) GetUserActiveAddress(ctx *fiber.Ctx) error {

	user, err := addressHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	result, err := addressHandler.AddressServices.GetUserAddressActive(user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get address",
		Data: map[string]any{
			"address": result,
		},
	})
}
