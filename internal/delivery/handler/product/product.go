package product

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandlerImpl struct {
	domain.ProductUseCase
	*validator.Validate
}

func NewProductHandler(product domain.ProductUseCase, v *validator.Validate) domain.ProductHandler {
	return &ProductHandlerImpl{
		ProductUseCase: product,
		Validate:       v,
	}
}

func (productHandler *ProductHandlerImpl) GetAllProduct(ctx *fiber.Ctx) error {

	result, err := productHandler.ProductUseCase.GetAllProduct()
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get products", map[string]any{
		"products": result,
	})
}

func (productHandler *ProductHandlerImpl) GetProductById(ctx *fiber.Ctx) error {

	paramProductId, err := ctx.ParamsInt("productId")
	if err != nil {
		return pkg.ErrCustomInvalidProductId
	}

	result, err := productHandler.ProductUseCase.GetProductById(int64(paramProductId))
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get product", map[string]any{
		"products": result,
	})

}

func (productHandler *ProductHandlerImpl) CreateProduct(ctx *fiber.Ctx) error {

	var body request.ReqCreateProduct

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := productHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := productHandler.ProductUseCase.CreateProduct(&body)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusCreated, "success create product", map[string]any{
		"products": result,
	})

}
func (productHandler *ProductHandlerImpl) UpdateProductById(ctx *fiber.Ctx) error {
	var body request.ReqUpdateProduct

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	productId, err := ctx.ParamsInt("productId")
	if err != nil {
		return pkg.ErrCustomInvalidProductId
	}

	if err = productHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := productHandler.ProductUseCase.UpdateProductById(&body, int64(productId))
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success update product", map[string]any{
		"products": result,
	})
}
