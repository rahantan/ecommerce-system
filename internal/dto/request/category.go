package request

type ReqCreateCategory struct {
	Name string `json:"name" validate:"required,max=20"`
}
type ReqUpdateCategory struct {
	ID   int64
	Name string `json:"name" validate:"required,max=20"`
}
