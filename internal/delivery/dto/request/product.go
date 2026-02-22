package request

type ReqCreateProduct struct {
	Name       string `json:"name" form:"name" validate:"required,max=20"`
	Price      int64  `json:"price" form:"price" validate:"required,gt=0"`
	Stock      int    `json:"stock" form:"stock" validate:"required,gte=0"`
	CategoryId int64  `json:"category_id" form:"category_id" validate:"required"`
}

type ReqUpdateProduct struct {
	Name       string `json:"name" form:"name" validate:"max=20"`
	Price      int64  `json:"price" form:"price" validate:"numeric"`
	Stock      int    `json:"stock" form:"stock" validate:"numeric"`
	CategoryId int64  `json:"category_id" form:"category_id" validate:"required,numeric"`
}
