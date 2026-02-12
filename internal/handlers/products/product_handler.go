package producthandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	productservices "ecommerce-system/internal/services/products"

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

func (productHandler *ProductHandlerImpl) GetAllProduct(ctx *fiber.Ctx) error {

	result, err := productHandler.ProductServices.GetAllProduct()
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get product",
		Data: map[string]any{
			"products": result,
		},
	})
}

func (productHandler *ProductHandlerImpl) GetProductById(ctx *fiber.Ctx) error {

	paramProductId, err := strconv.Atoi(ctx.Params("productId"))
	if err != nil {
		return exceptions.ErrCustomInvalidProductId
	}

	result, err := productHandler.ProductServices.GetProductById(int64(paramProductId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get product",
		Data: map[string]any{
			"product": result,
		},
	})
}

func (productHandler *ProductHandlerImpl) CreateProduct(ctx *fiber.Ctx) error {

	var body request.ReqCreateProduct

	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := productHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	result, err := productHandler.ProductServices.CreateProduct(&body)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success create product",
		Data: map[string]any{
			"product": result,
		},
	})

}
func (productHandler *ProductHandlerImpl) UpdateProductById(ctx *fiber.Ctx) error {
	var body request.ReqUpdateProduct

	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	productId, err := strconv.Atoi(ctx.Params("productId"))
	if err != nil {
		return exceptions.ErrCustomInvalidProductId
	}

	if err = productHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	result, err := productHandler.ProductServices.UpdateProductById(&body, int64(productId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success update product",
		Data: map[string]any{
			"product": result,
		},
	})
}
