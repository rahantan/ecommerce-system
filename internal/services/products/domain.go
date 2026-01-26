package productservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type ProductServices interface {
	GetProductById(id int64) (*response.ResProduct, error)
	GetAllProduct() ([]*response.ResProduct, error)
	CreateProduct(request *request.ReqCreateProduct) (*response.ResProduct, error)
	UpdateProduct(request *request.ReqUpdateProduct) (*response.ResProduct, error)
}
