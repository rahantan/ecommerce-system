package productrepositories

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	*gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepositories {
	return &ProductRepositoryImpl{
		DB: db,
	}
}

func (productRepo *ProductRepositoryImpl) GetProductById(id int64) (*models.ProductModel, error) {
	var product models.ProductModel
	err := productRepo.DB.Preload("Category").Where("id=?", id).Take(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}
func (productRepo *ProductRepositoryImpl) CreateProduct(product *models.ProductModel) (*models.ProductModel, error) {

	if err := productRepo.DB.Create(product).Error; err != nil {
		if containErr := exceptions.CheckContainError(err); containErr {
			return nil, exceptions.ErrCategoryNotFound
		}
		return nil, err
	}

	return productRepo.GetProductById(product.ID)
}
func (productRepo *ProductRepositoryImpl) UpdateProductById(product *models.ProductModel) (*models.ProductModel, error) {
	result := productRepo.DB.Model(&models.ProductModel{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, exceptions.ErrNoRowsAffected
	}
	return productRepo.GetProductById(product.ID)
}

func (productRepo *ProductRepositoryImpl) GetAllProduct() ([]*models.ProductModel, error) {
	var products []*models.ProductModel
	err := productRepo.DB.Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
