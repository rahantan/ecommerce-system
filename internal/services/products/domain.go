package productservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type ProductServices interface {
	GetProductById(productID int64) (*response.ResProduct, error)
	GetAllProduct() ([]*response.ResProduct, error)
	CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error)
	UpdateProductById(request *request.ReqUpdateProduct, productID int64) (*response.ResProduct, error)
}
