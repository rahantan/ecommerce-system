package productrepositories

import "ecommerce-system/internal/models"

type ProductRepositories interface {
	GetAllProduct() ([]*models.ProductModel, error)
	UpdateProductById(product *models.ProductModel) (*models.ProductModel, error)
	CreateProduct(product *models.ProductModel) (*models.ProductModel, error)
	GetProductById(id int64) (*models.ProductModel, error)
}
