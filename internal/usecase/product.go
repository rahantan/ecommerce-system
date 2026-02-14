package usecase

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

type ProductUseCaseImpl struct {
	domain.ProductRepository
	*gorm.DB
}

func NewProductUseCase(product domain.ProductRepository, db *gorm.DB) domain.ProductUseCase {
	return &ProductUseCaseImpl{
		ProductRepository: product,
		DB:                db,
	}
}

func (product *ProductUseCaseImpl) loadProduct(productLoad *model.ProductModel) *response.ResProduct {
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
func (productUC *ProductUseCaseImpl) GetProductById(productID int64) (*response.ResProduct, error) {
	result, err := productUC.ProductRepository.GetProductById(productUC.DB, productID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return productUC.loadProduct(result), nil
}
func (productUC *ProductUseCaseImpl) GetAllProduct() ([]*response.ResProduct, error) {

	result, err := productUC.ProductRepository.GetAllProduct(productUC.DB)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	products := []*response.ResProduct{}
	for _, product := range result {
		products = append(products, productUC.loadProduct(product))
	}

	return products, nil
}

func (productUC *ProductUseCaseImpl) CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error) {

	result, err := productUC.ProductRepository.CreateProduct(productUC.DB, &model.ProductModel{
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})

	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return productUC.loadProduct(result), nil
}
func (productUC *ProductUseCaseImpl) UpdateProductById(request *request.ReqUpdateProduct, productID int64) (*response.ResProduct, error) {
	if err := productUC.ProductRepository.CheckProductNotFoundForUpdate(productUC.DB, productID); err != nil {
		return nil, pkg.MappingError(err)
	}

	result, err := productUC.ProductRepository.UpdateProductById(productUC.DB, &model.ProductModel{
		ID:         productID,
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
	})

	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return productUC.loadProduct(result), nil
}
