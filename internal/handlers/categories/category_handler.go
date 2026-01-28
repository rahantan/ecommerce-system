package categoryhandlers

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	categoryservices "ecommerce-system/internal/services/categories"
	"ecommerce-system/internal/utils"
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

func (categoryHandler *CategoryHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
}
func (categoryHandler *CategoryHandlerImpl) CreateCategory(ctx *fiber.Ctx) error {
	var body request.ReqCreateCategory

	err := ctx.BodyParser(&body)
	if err != nil {
		return categoryHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to create category")
	}

	err = categoryHandler.Validate.Struct(&body)
	if err != nil {
		return categoryHandler.withMessage(exceptions.ValidationError(err), "failed to create category")
	}

	result, err := categoryHandler.CategoryServices.CreateCategory(&body)
	if err != nil {
		return categoryHandler.withMessage(err, "failed to create category")
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
		return categoryHandler.withMessage(err, "failed to get categories")
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

	err := ctx.BodyParser(&body)
	if err != nil {
		return categoryHandler.withMessage(exceptions.ErrCustomInvalidPayload, "failed to update category")
	}

	err = categoryHandler.Validate.Struct(&body)
	if err != nil {
		return categoryHandler.withMessage(exceptions.ValidationError(err), "failed to update category")
	}

	result, err := categoryHandler.CategoryServices.UpdateCategory(&body)
	if err != nil {
		return categoryHandler.withMessage(err, "failed to update category")
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
		return categoryHandler.withMessage(exceptions.ErrCustomInvalidCategoryId, "failed to get category")
	}

	result, err := categoryHandler.CategoryServices.GetCategoryById(int64(paramCategId))
	if err != nil {
		return categoryHandler.withMessage(err, "failed to get category")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get category",
		Data: map[string]any{
			"category": result,
		},
	})
}
