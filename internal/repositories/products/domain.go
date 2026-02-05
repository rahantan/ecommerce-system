package productrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type ProductRepositories interface {
	GetAllProduct(db *gorm.DB) ([]*models.ProductModel, error)
	GetAllProductByIDs(db *gorm.DB, productIDs []int64) ([]*models.ProductModel, error)
	UpdateProductById(db *gorm.DB, product *models.ProductModel) (*models.ProductModel, error)
	CreateProduct(db *gorm.DB, product *models.ProductModel) (*models.ProductModel, error)
	GetProductById(db *gorm.DB, id int64) (*models.ProductModel, error)
	CheckProductNotFoundForUpdate(db *gorm.DB, productID int64) error

	//utils
	UpdateProductStockByID(db *gorm.DB, product []*models.ProductModel) error
}
