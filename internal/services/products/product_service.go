package productservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	productrepositories "ecommerce-system/internal/repositories/products"
)

type ProductServiceImpl struct {
	productrepositories.ProductRepositories
}

func NewProductService(product productrepositories.ProductRepositories) ProductServices {
	return &ProductServiceImpl{
		ProductRepositories: product,
	}
}
func (productService *ProductServiceImpl) handleError(err error) error {
	return exceptions.CheckError(err)
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
func (productService *ProductServiceImpl) GetProductById(id int64) (*response.ResProduct, error) {
	result, err := productService.ProductRepositories.GetProductById(id)
	if errCheck := productService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return productService.loadProduct(result), nil
}
func (productService *ProductServiceImpl) GetAllProduct() ([]*response.ResProduct, error) {
	result, err := productService.ProductRepositories.GetAllProduct()
	if errCheck := productService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	products := []*response.ResProduct{}
	for _, product := range result {
		products = append(products, productService.loadProduct(product))
	}

	return products, nil
}

func (productService *ProductServiceImpl) CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error) {

	result, err := productService.ProductRepositories.CreateProduct(&models.ProductModel{
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})

	if errCheck := productService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	return productService.loadProduct(result), nil
}
func (productService *ProductServiceImpl) UpdateProduct(request *request.ReqUpdateProduct) (*response.ResProduct, error) {
	result, err := productService.ProductRepositories.UpdateProductById(&models.ProductModel{
		ID:         request.ID,
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})
	if errCheck := productService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return productService.loadProduct(result), nil
}
