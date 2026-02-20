package request

type ReqCreateOrUpdateCartItem struct {
	ProductID int64 `json:"product_id" validate:"required,numeric"`
	Qty       int   `json:"qty" validate:"required,numeric,gt=0"`
}

type ReqUpdateCartQty struct {
	// CartID int64 `json:"cart_id" validate:"required,numeric"`
	Qty int `json:"qty" validate:"required,numeric,gt=0"`
}

type ReqAddCart struct {
	ProductID int64 `json:"product_id" validate:"required,numeric"`
}

type ReqDeleteCartItem struct {
	CartIDs []int64 `json:"cart_ids" validate:"required,min=1,dive,required"`
}
