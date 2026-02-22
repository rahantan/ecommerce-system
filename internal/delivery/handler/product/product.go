package product

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"
	"os"
	"strconv"

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

func (productHandler *ProductHandlerImpl) GetProductImage(ctx *fiber.Ctx) error {
	fileName := ctx.Query("src")
	if fileName == "" {
		return pkg.NewError(pkg.KindBadRequest, "invalid image name", nil)
	}

	filePath := pkg.ProductDir + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return pkg.NewError(pkg.KindNotFound, "image not found", nil)
	}

	return ctx.SendFile(pkg.ProductDir + fileName)
}
func (productHandler *ProductHandlerImpl) GetAllProduct(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return pkg.ErrInvalidPage
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil {
		return pkg.ErrinvalidLimit
	}

	result, meta, err := productHandler.ProductUseCase.GetAllProduct(page, limit)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get products", map[string]any{
		"products": result,
		"meta":     meta,
	})
}

func (productHandler *ProductHandlerImpl) GetProductById(ctx *fiber.Ctx) error {

	productID, err := ctx.ParamsInt("productId")
	if err != nil {
		return pkg.ErrCustomInvalidProductId
	}

	result, err := productHandler.ProductUseCase.GetProductById(int64(productID))
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get product", map[string]any{
		"products": result,
	})

}

func (productHandler *ProductHandlerImpl) DeleteProductByID(ctx *fiber.Ctx) error {
	productID, err := ctx.ParamsInt("productId")
	if err != nil {
		return pkg.ErrCustomInvalidProductId
	}
	if err := productHandler.ProductUseCase.DeleteProductById(int64(productID)); err != nil {
		return err
	}
	return handler.SuccessResponse(ctx, fiber.StatusOK, "success delete product", nil)
}

func (productHandler *ProductHandlerImpl) CreateProduct(ctx *fiber.Ctx) error {

	var body request.ReqCreateProduct

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := productHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	productImage := &pkg.File{
		Reader:   file,
		Filename: fileHeader.Filename,
		Src:      pkg.ProductDir,
	}

	result, err := productHandler.ProductUseCase.CreateProduct(&body, productImage)
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

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	productImage := &pkg.File{
		Reader:   file,
		Filename: fileHeader.Filename,
		Src:      pkg.ProductDir,
	}

	result, err := productHandler.ProductUseCase.UpdateProductById(&body, productImage, int64(productId))
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success update product", map[string]any{
		"products": result,
	})
}
