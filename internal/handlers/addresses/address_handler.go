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
func (addressHandler *AddressHandlerImpl) CreateAddress(ctx *fiber.Ctx) error {
	var body request.ReqCreateAddress

	user := ctx.Locals("user").(response.ResUser)
	body.UserID = user.ID

	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.UpdateMessageErr(exceptions.ErrCustomInvalidPayload, exceptions.MsgFailCreateAddress)
	}

	err = addressHandler.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailCreateAddress)
	}

	result, err := addressHandler.AddressServices.CreateAddress(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailCreateAddress)
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

	user := ctx.Locals("user").(response.ResUser)

	result, err := addressHandler.AddressServices.GetAllAddress(user.ID)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailGetAllAddresses)
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
		return utils.UpdateMessageErr(exceptions.ErrCustomInvalidPayload, exceptions.MsgFailUpdateAddress)
	}

	addressId, err := strconv.Atoi(ctx.Params("addressId"))
	if err != nil {
		return utils.UpdateMessageErr(exceptions.ErrCustomInvalidAddressId, exceptions.MsgFailUpdateAddress)
	}

	user := ctx.Locals("user").(response.ResUser)
	body.UserID = user.ID
	body.ID = int64(addressId)

	err = addressHandler.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailUpdateAddress)
	}

	result, err := addressHandler.AddressServices.UpdateAddressByUserId(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailUpdateAddress)
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
	user := ctx.Locals("user").(response.ResUser)

	result, err := addressHandler.AddressServices.GetUserActiveAddress(user.ID)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailGetAddress)
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get address",
		Data: map[string]any{
			"address": result,
		},
	})
}
