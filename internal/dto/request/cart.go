package request

type ReqCreateOrUpdateCartItem struct {
	ProductID int64 `json:"product_id" validate:"required,numeric"`
	Qty       int   `json:"qty" validate:"required,numeric,gt=0"`
}

type ReqDeleteCartItem struct {
	CartIDs []int64 `json:"cart_ids" validate:"required,min=1,dive,required"`
}

// type ReqUpdateCartItem struct {
// 	ProductID int64 `json:"product_id" validate:"required,numeric"`
// 	Qty       int64 `json:"qty" validate:"required,numeric,gt=0"`
// }
