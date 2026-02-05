package addresshandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	addressservices "ecommerce-system/internal/services/addresses"
	"ecommerce-system/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AddressHandlerImpl struct {
	addressservices.AddressServices
	*validator.Validate
}

func NewAddressHandler(address addressservices.AddressServices, v *validator.Validate) AddressHandlers {
	return &AddressHandlerImpl{
		AddressServices: address,
		Validate:        v,
	}
}

func (addressHandler *AddressHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
}
func (addressHandler *AddressHandlerImpl) CreateAddress(ctx *fiber.Ctx) error {
	var body request.ReqCreateAddress

	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return addressHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to create address")
	}
	err := ctx.BodyParser(&body)
	if err != nil {
		return addressHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to create address")
	}

	err = addressHandler.Validate.Struct(&body)
	if err != nil {
		return addressHandler.withMessage(exceptions.ValidationError(err), "failed to create address")
	}

	result, err := addressHandler.AddressServices.CreateAddress(&body, user.ID)
	if err != nil {
		return addressHandler.withMessage(err, "failed to create address")
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

	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return addressHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to get addresses")
	}
	result, err := addressHandler.AddressServices.GetAllAddress(user.ID)
	if err != nil {
		return addressHandler.withMessage(err, "failed to get addresses")
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
	var body request.ReqUpdateAddress

	err := ctx.BodyParser(&body)
	if err != nil {
		return addressHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to update address")
	}

	addressId, err := strconv.Atoi(ctx.Params("addressId"))
	if err != nil {
		return addressHandler.withMessage(exceptions.ErrCustomInvalidAddressId, "failed to update address")
	}

	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return addressHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to update address")
	}
	err = addressHandler.Validate.Struct(&body)
	if err != nil {
		return addressHandler.withMessage(exceptions.ValidationError(err), "failed to update address")
	}

	result, err := addressHandler.AddressServices.UpdateAddressByUserId(&body, int64(addressId), user.ID)
	if err != nil {
		return addressHandler.withMessage(err, "failed to update address")
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

	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return addressHandler.withMessage(exceptions.ErrCustomUnauthorized, "failed to get address")
	}
	result, err := addressHandler.AddressServices.GetUserAddressActive(user.ID)
	if err != nil {
		return addressHandler.withMessage(err, "failed to get address")
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get address",
		Data: map[string]any{
			"address": result,
		},
	})
}
