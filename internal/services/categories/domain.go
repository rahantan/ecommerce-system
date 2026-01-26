package categoryservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type CategoryServices interface {
	GetCategoryById(id int64) (*response.ResCategory, error)

	GetAllCategory() ([]*response.ResCategory, error)
	CreateCategory(request *request.ReqCreateCategory) (*response.ResCategory, error)
	UpdateCategory(request *request.ReqUpdateCategory) (*response.ResCategory, error)
}
