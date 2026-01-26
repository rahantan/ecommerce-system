package producthandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	productservices "ecommerce-system/internal/services/products"
	"ecommerce-system/internal/utils"
	"fmt"

	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandlerImpl struct {
	productservices.ProductServices
	*validator.Validate
}

func NewProductHandler(product productservices.ProductServices, v *validator.Validate) ProductHandlers {
	return &ProductHandlerImpl{
		ProductServices: product,
		Validate:        v,
	}
}
func (productHandler *ProductHandlerImpl) CreateProduct(ctx *fiber.Ctx) error {
	fmt.Println("4")
	var body request.ReqCreateProduct

	err := ctx.BodyParser(&body)
	if err != nil {
		return exceptions.ErrCustomInvalidPayload.WithMessage(exceptions.MsgFailCreateProduct)
	}

	err = productHandler.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailCreateProduct)
	}

	result, err := productHandler.ProductServices.CreateProduct(&body)
	if err != nil {

		return utils.UpdateMessageErr(err, exceptions.MsgFailCreateProduct)
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success create product",
		Data: map[string]any{
			"product": result,
		},
	})

}
func (productHandler *ProductHandlerImpl) GetAllProduct(ctx *fiber.Ctx) error {

	result, err := productHandler.ProductServices.GetAllProduct()
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailGetAllProduct)
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get product",
		Data: map[string]any{
			"products": result,
		},
	})
}
func (productHandler *ProductHandlerImpl) UpdateProductById(ctx *fiber.Ctx) error {
	var body request.ReqUpdateProduct

	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailUpdateProduct)
	}

	err = productHandler.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailUpdateProduct)
	}

	result, err := productHandler.ProductServices.UpdateProduct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailUpdateProduct)
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success update product",
		Data: map[string]any{
			"product": result,
		},
	})
}
func (productHandler *ProductHandlerImpl) GetProductById(ctx *fiber.Ctx) error {
	paramProductId, err := strconv.Atoi(ctx.Params("productId"))
	if err != nil {
		return exceptions.ErrCustomInvalidProductId.WithMessage(exceptions.MsgFailGetProduct)
	}

	result, err := productHandler.ProductServices.GetProductById(int64(paramProductId))
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailGetProduct)
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get product",
		Data: map[string]any{
			"product": result,
		},
	})
}
