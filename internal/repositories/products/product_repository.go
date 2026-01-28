package productrepositories

import (
	"ecommerce-system/internal/models"

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

func (productRepo *ProductRepositoryImpl) checkErrMysql(err error) error {
	if models.IsInternalErrMysql(err) {
		return err
	}
	if models.ForeignKeyErr(err) {
		return models.ErrCategoryNotFound
	}
	return models.ErrProductNotFound
}
func (productRepo *ProductRepositoryImpl) checkProductNotFoundForUpdate(productID int64) bool {
	var count int64

	if err := productRepo.DB.Model(&models.ProductModel{}).Where("id=?", productID).Count(&count).Error; err != nil {
		return false
	}

	return count == 0
}

func (productRepo *ProductRepositoryImpl) GetProductById(productID int64) (*models.ProductModel, error) {
	var product models.ProductModel

	if err := productRepo.DB.Preload("Category").Where("id=?", productID).Take(&product).Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}

	return &product, nil
}

func (productRepo *ProductRepositoryImpl) CreateProduct(product *models.ProductModel) (*models.ProductModel, error) {
	if err := productRepo.DB.Create(product).Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}
	return productRepo.GetProductById(product.ID)
}

func (productRepo *ProductRepositoryImpl) UpdateProductById(product *models.ProductModel) (*models.ProductModel, error) {

	if productRepo.checkProductNotFoundForUpdate(product.ID) {
		return nil, models.ErrProductNotFound
	}

	result := productRepo.DB.Model(&models.ProductModel{}).Where("id=?", product.ID).Updates(product)
	if err := result.Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}

	if result.RowsAffected < 1 {
		return nil, models.ErrNoRowsAffected
	}

	return productRepo.GetProductById(product.ID)
}

func (productRepo *ProductRepositoryImpl) GetAllProduct() ([]*models.ProductModel, error) {
	var products []*models.ProductModel

	if err := productRepo.DB.Preload("Category").Find(&products).Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}
	return products, nil
}
