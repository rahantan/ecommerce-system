package request

type ReqCreateProduct struct {
	Name       string  `json:"name" validate:"required,max=20"`
	Price      float64 `json:"price" validate:"required,numeric"`
	Stock      int     `json:"stock" validate:"required,numeric"`
	CategoryId int64   `json:"category_id" validate:"required,numeric"`
}
type ReqUpdateProduct struct {
	ID         int64   `json:"id" validate:"required,numeric"`
	Name       string  `json:"name" validate:"max=20"`
	Price      float64 `json:"price" validate:"numeric"`
	Stock      int     `json:"stock" validate:"numeric"`
	CategoryId int64   `json:"category_id" validate:"numeric"`
}
