package request

type ReqCreateProduct struct {
	Name       string `json:"name" validate:"required,max=20"`
	Price      int64  `json:"price" validate:"required,numeric,gt=0"`
	Stock      int    `json:"stock" validate:"required,numeric,gte=0"`
	CategoryId int64  `json:"category_id" validate:"required,numeric"`
}
type ReqUpdateProduct struct {
	Name       string `json:"name" validate:"max=20"`
	Price      int64  `json:"price" validate:"numeric,gt=0"`
	Stock      int    `json:"stock" validate:"numeric,gte=0"`
	CategoryId int64  `json:"category_id" validate:"required,numeric"`
}
