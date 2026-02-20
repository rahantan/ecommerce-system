package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAllProduct(db *gorm.DB) ([]*model.ProductModel, error)
	FindByIDs(db *gorm.DB, productIDs []int64) (map[int64]*model.ProductModel, error)
	FindByIDsForUpdate(db *gorm.DB, productIDs []int64) (map[int64]*model.ProductModel, error)
	UpdateProductById(db *gorm.DB, product *model.ProductModel) (*model.ProductModel, error)
	CreateProduct(db *gorm.DB, product *model.ProductModel) (*model.ProductModel, error)
	GetProductById(db *gorm.DB, id int64) (*model.ProductModel, error)
	CheckProductNotFoundForUpdate(db *gorm.DB, productID int64) error

	// FindProductsByIDs(db *gorm.DB, ids []int64) (map[int64]*model.ProductModel, error)

	UpdateProductStock(db *gorm.DB, products []model.ProductModel) error
}

type ProductUseCase interface {
	GetProductById(productID int64) (*response.ResProduct, error)
	GetAllProduct() ([]*response.ResProduct, error)
	CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error)
	UpdateProductById(request *request.ReqUpdateProduct, productID int64) (*response.ResProduct, error)
}

type ProductHandler interface {
	CreateProduct(ctx *fiber.Ctx) error
	GetAllProduct(ctx *fiber.Ctx) error
	UpdateProductById(ctx *fiber.Ctx) error
	GetProductById(ctx *fiber.Ctx) error
}
