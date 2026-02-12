package categoryhandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	categoryservices "ecommerce-system/internal/services/categories"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandlerImpl struct {
	categoryservices.CategoryServices
	*validator.Validate
}

func NewCategoryHandler(category categoryservices.CategoryServices, v *validator.Validate) CategoryHandlers {
	return &CategoryHandlerImpl{
		CategoryServices: category,
		Validate:         v,
	}
}

func (categoryHandler *CategoryHandlerImpl) CreateCategory(ctx *fiber.Ctx) error {
	var body request.ReqCreateCategory

	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	if err := categoryHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	result, err := categoryHandler.CategoryServices.CreateCategory(&body)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success create category",
		Data: map[string]any{
			"category": result,
		},
	})
}
func (categoryHandler *CategoryHandlerImpl) GetAllCategory(ctx *fiber.Ctx) error {

	result, err := categoryHandler.CategoryServices.GetAllCategory()
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get category",
		Data: map[string]any{
			"categories": result,
		},
	})
}
func (categoryHandler *CategoryHandlerImpl) UpdateCategoryById(ctx *fiber.Ctx) error {

	var body request.ReqUpdateCategory

	if err := ctx.BodyParser(&body); err != nil {
		return exceptions.ErrCustomInvalidPayload
	}

	categoryId, err := strconv.Atoi(ctx.Params("productId"))
	if err != nil {
		return exceptions.ErrCustomInvalidCategoryId
	}

	if err = categoryHandler.Validate.Struct(&body); err != nil {
		return exceptions.ValidationError(err)
	}

	result, err := categoryHandler.CategoryServices.UpdateCategory(&body, int64(categoryId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success update category",
		Data: map[string]any{
			"category": result,
		},
	})
}
func (categoryHandler *CategoryHandlerImpl) GetCategoryById(ctx *fiber.Ctx) error {

	paramCategId, err := strconv.Atoi(ctx.Params("categoryId"))
	if err != nil {
		return exceptions.ErrCustomInvalidCategoryId
	}

	result, err := categoryHandler.CategoryServices.GetCategoryById(int64(paramCategId))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get category",
		Data: map[string]any{
			"category": result,
		},
	})
}
