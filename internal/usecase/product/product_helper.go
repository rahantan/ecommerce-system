package product

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
)

func (product *ProductUseCaseImpl) resProduct(productLoad *model.ProductModel) *response.ResProduct {
	return &response.ResProduct{
		ID:        productLoad.ID,
		Name:      productLoad.Name,
		Price:     productLoad.Price,
		Stock:     productLoad.Stock,
		ImageUrl:  "localhost:8080/api/public/products/image?src=" + productLoad.Image,
		CreatedAt: productLoad.CreatedAt.Format(pkg.DateTimeLayout),
		UpdatedAt: productLoad.UpdatedAt.Format(pkg.DateTimeLayout),
		ResCategory: &response.ResCategory{
			ID:   productLoad.Category.ID,
			Name: productLoad.Category.Name,
		},
	}
}
