package cartitemsservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type CartItemServices interface {
	CreateOrUpdateCartItem(request *request.ReqCreateOrUpdateCartItem, userID int64) (*response.ResCartItem, error)
	GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error)
	DeleteCartItemById(cartID int64, userID int64) error
	GetCartItemByProductUser(productID int64, userID int64) (*response.ResCartItem, error)
}
