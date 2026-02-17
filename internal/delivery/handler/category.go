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

type CategoryHandlerImpl struct {
	domain.CategoryUseCase
	*validator.Validate
}

func NewCategoryHandler(category domain.CategoryUseCase, v *validator.Validate) domain.CategoryHandler {
	return &CategoryHandlerImpl{
		CategoryUseCase: category,
		Validate:        v,
	}
}

func (categoryHandler *CategoryHandlerImpl) CreateCategory(ctx *fiber.Ctx) error {
	var body request.ReqCreateCategory

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := categoryHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := categoryHandler.CategoryUseCase.CreateCategory(&body)
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

	result, err := categoryHandler.CategoryUseCase.GetAllCategory()
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
		return pkg.ErrCustomInvalidPayload
	}

	categoryId, err := strconv.Atoi(ctx.Params("categoryId"))
	if err != nil {
		return pkg.ErrCustomInvalidCategoryId
	}

	if err = categoryHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := categoryHandler.CategoryUseCase.UpdateCategory(&body, int64(categoryId))
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
		return pkg.ErrCustomInvalidCategoryId
	}

	result, err := categoryHandler.CategoryUseCase.GetCategoryById(int64(paramCategId))
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
