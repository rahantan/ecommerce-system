package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoryById(db *gorm.DB, categoryID int64) (*model.CategoryModel, error)
	GetAllCategory(db *gorm.DB) ([]*model.CategoryModel, error)
	UpdateCategoryById(db *gorm.DB, category *model.CategoryModel) (*model.CategoryModel, error)
	CreateCategory(db *gorm.DB, category *model.CategoryModel) (*model.CategoryModel, error)
}

type CategoryUseCase interface {
	GetCategoryById(id int64) (*response.ResCategory, error)
	GetAllCategory() ([]*response.ResCategory, error)
	CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error)
	UpdateCategory(request *request.ReqUpdateCategory, categoryID int64) (*response.ResCategory, error)
}

type CategoryHandler interface {
	CreateCategory(ctx *fiber.Ctx) error
	GetAllCategory(ctx *fiber.Ctx) error
	UpdateCategoryById(ctx *fiber.Ctx) error
	GetCategoryById(ctx *fiber.Ctx) error
}
