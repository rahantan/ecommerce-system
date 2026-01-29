package producthandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	productservices "ecommerce-system/internal/services/products"
	"ecommerce-system/internal/utils"

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

func (productHandler *ProductHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
}

func (productHandler *ProductHandlerImpl) CreateProduct(ctx *fiber.Ctx) error {

	var body request.ReqCreateProduct

	err := ctx.BodyParser(&body)
	if err != nil {
		return productHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to create product")
	}

	err = productHandler.Validate.Struct(&body)
	if err != nil {
		return productHandler.withMessage(exceptions.ValidationError(err), "failed to create product")
	}

	result, err := productHandler.ProductServices.CreateProduct(&body)
	if err != nil {
		return productHandler.withMessage(err, "failed to create product")
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
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
		return productHandler.withMessage(err, "failed to get products")
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
		return productHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to update product")
	}

	productId, err := strconv.Atoi(ctx.Params("productId"))
	if err != nil {
		return productHandler.withMessage(exceptions.ErrCustomInvalidProductId, "failed to update address")
	}
	body.ID = int64(productId)
	err = productHandler.Validate.Struct(&body)
	if err != nil {
		return productHandler.withMessage(exceptions.ValidationError(err), "failed to update product")
	}

	result, err := productHandler.ProductServices.UpdateProduct(&body)
	if err != nil {
		return productHandler.withMessage(err, "failed to update product")
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
		return productHandler.withMessage(exceptions.ErrCustomInvalidProductId, "failed to get product")
	}

	result, err := productHandler.ProductServices.GetProductById(int64(paramProductId))
	if err != nil {
		return productHandler.withMessage(err, "failed to get product")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get product",
		Data: map[string]any{
			"product": result,
		},
	})
}
