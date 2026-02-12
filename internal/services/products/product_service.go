package productservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/models"
	productrepositories "ecommerce-system/internal/repositories/products"
	"ecommerce-system/internal/utils"

	"gorm.io/gorm"
)

type ProductServiceImpl struct {
	productrepositories.ProductRepositories
	*gorm.DB
}

func NewProductService(product productrepositories.ProductRepositories, db *gorm.DB) ProductServices {
	return &ProductServiceImpl{
		ProductRepositories: product,
		DB:                  db,
	}
}

func (product *ProductServiceImpl) loadProduct(productLoad *models.ProductModel) *response.ResProduct {
	return &response.ResProduct{
		ID:        productLoad.ID,
		Name:      productLoad.Name,
		Price:     productLoad.Price,
		Stock:     productLoad.Stock,
		CreatedAt: productLoad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: productLoad.UpdatedAt.Format("2006-01-02 15:04:05"),
		ResCategory: &response.ResCategory{
			ID:   productLoad.Category.ID,
			Name: productLoad.Category.Name,
		},
	}
}
func (productService *ProductServiceImpl) GetProductById(productID int64) (*response.ResProduct, error) {
	result, err := productService.ProductRepositories.GetProductById(productService.DB, productID)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	return productService.loadProduct(result), nil
}
func (productService *ProductServiceImpl) GetAllProduct() ([]*response.ResProduct, error) {

	result, err := productService.ProductRepositories.GetAllProduct(productService.DB)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	products := []*response.ResProduct{}
	for _, product := range result {
		products = append(products, productService.loadProduct(product))
	}

	return products, nil
}

func (productService *ProductServiceImpl) CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error) {

	result, err := productService.ProductRepositories.CreateProduct(productService.DB, &models.ProductModel{
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})

	if err != nil {
		return nil, utils.MappingError(err)
	}

	return productService.loadProduct(result), nil
}
func (productService *ProductServiceImpl) UpdateProductById(request *request.ReqUpdateProduct, productID int64) (*response.ResProduct, error) {
	if err := productService.ProductRepositories.CheckProductNotFoundForUpdate(productService.DB, productID); err != nil {
		return nil, utils.MappingError(err)
	}

	result, err := productService.ProductRepositories.UpdateProductById(productService.DB, &models.ProductModel{
		ID:         productID,
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})

	if err != nil {
		return nil, utils.MappingError(err)
	}

	return productService.loadProduct(result), nil
}
