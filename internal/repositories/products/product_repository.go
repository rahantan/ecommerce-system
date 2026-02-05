package productrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository(db *gorm.DB) ProductRepositories {
	return &ProductRepositoryImpl{}
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
func (productRepo *ProductRepositoryImpl) CheckProductNotFoundForUpdate(db *gorm.DB, productID int64) error {
	var count int64
	if err := db.Model(&models.ProductModel{}).Where("id=?", productID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return models.ErrProductNotFound
	}
	return nil
}
func (productRepo *ProductRepositoryImpl) GetAllProductByIDs(db *gorm.DB, productIDs []int64) ([]*models.ProductModel, error) {
	var products []*models.ProductModel

	if err := db.Where("id IN ?", productIDs).Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (productRepo *ProductRepositoryImpl) GetProductById(db *gorm.DB, productID int64) (*models.ProductModel, error) {
	var product models.ProductModel

	if err := db.Preload("Category").Where("id=?", productID).Take(&product).Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}

	return &product, nil
}

func (productRepo *ProductRepositoryImpl) CreateProduct(db *gorm.DB, product *models.ProductModel) (*models.ProductModel, error) {
	if err := db.Create(product).Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}
	return productRepo.GetProductById(db, product.ID)
}

func (productRepo *ProductRepositoryImpl) UpdateProductById(db *gorm.DB, product *models.ProductModel) (*models.ProductModel, error) {

	result := db.Model(&models.ProductModel{}).Where("id=?", product.ID).Updates(product)
	if err := result.Error; err != nil {
		return nil, productRepo.checkErrMysql(err)
	}

	if result.RowsAffected < 1 {
		return nil, models.ErrNoRowsAffected
	}

	return productRepo.GetProductById(db, product.ID)
}

func (productRepo *ProductRepositoryImpl) GetAllProduct(db *gorm.DB) ([]*models.ProductModel, error) {
	var products []*models.ProductModel

	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (productRepo *ProductRepositoryImpl) UpdateProductStockByID(db *gorm.DB, products []*models.ProductModel) error {

	sql := "UPDATE products SET stock = CASE id "
	args := []interface{}{}
	ids := []interface{}{}

	for _, p := range products {
		sql += "WHEN ? THEN ? "
		args = append(args, p.ID, p.Stock)
		ids = append(ids, p.ID)
	}

	args = append(args, ids)

	sql += "END WHERE id IN (?)"
	return db.Exec(sql, args...).Error

}
