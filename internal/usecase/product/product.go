package product

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"math"

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

func (productUC *ProductUseCaseImpl) GetProductById(productID int64) (*response.ResProduct, error) {

	result, err := productUC.ProductRepository.GetProductById(productUC.DB, productID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return productUC.resProduct(result), nil
}

func (productUC *ProductUseCaseImpl) GetAllProduct(page, limit int) ([]*response.ResProduct, *response.ResPaginateStandard, error) {

	productsResulut, count, err := productUC.ProductRepository.GetAllProduct(productUC.DB, page, limit)
	if err != nil {
		return nil, nil, pkg.MappingError(err)
	}

	totalPage := math.Ceil(float64(count) / float64(limit))
	products := []*response.ResProduct{}
	for _, product := range productsResulut {
		products = append(products, productUC.resProduct(product))
	}

	paginate := &response.ResPaginateStandard{
		Page:      page,
		Limit:     limit,
		TotalData: int(count),
		TotalPage: int(totalPage),
	}
	return products, paginate, nil
}

func (productUC *ProductUseCaseImpl) CreateProduct(request *request.ReqCreateProduct, productImage *pkg.File) (*response.ResProduct, error) {

	if err := pkg.UploadFile(productImage); err != nil {
		return nil, pkg.MappingError(err)
	}

	result, err := productUC.ProductRepository.CreateProduct(productUC.DB, &model.ProductModel{
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
		Image:      productImage.Filename,
	})

	if err != nil {
		if err := pkg.Delete(pkg.ProductDir); err != nil {
			return nil, pkg.MappingError(err)
		}
		return nil, pkg.MappingError(err)
	}

	return productUC.resProduct(result), nil
}

func (productUC *ProductUseCaseImpl) UpdateProductById(request *request.ReqUpdateProduct, productImage *pkg.File, productID int64) (*response.ResProduct, error) {
	if err := pkg.UploadFile(productImage); err != nil {
		return nil, pkg.MappingError(err)
	}

	oldProduct, err := productUC.ProductRepository.GetProductById(productUC.DB, productID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	result, err := productUC.ProductRepository.UpdateProductById(productUC.DB, &model.ProductModel{
		ID:         productID,
		Name:       request.Name,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryID: request.CategoryId,
		Image:      productImage.Filename,
	})

	if err != nil {
		if err := pkg.Delete(pkg.ProductDir + result.Image); err != nil {
			return nil, pkg.MappingError(err)
		}
		return nil, pkg.MappingError(err)
	}

	if err := pkg.Delete(pkg.ProductDir + oldProduct.Image); err != nil {
		return nil, pkg.MappingError(err)
	}

	return productUC.resProduct(result), nil
}

func (productUC *ProductUseCaseImpl) DeleteProductById(productID int64) error {
	product, err := productUC.ProductRepository.GetProductById(productUC.DB, productID)
	if err != nil {
		return pkg.MappingError(err)
	}

	if err := productUC.ProductRepository.DeleteByID(productUC.DB, productID); err != nil {
		return pkg.MappingError(err)
	}

	if err := pkg.Delete(pkg.ProductDir + product.Image); err != nil {
		return err
	}

	return nil
}
