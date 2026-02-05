package orderservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

// type OrderServices interface {
// 	CreateOrder(req []*request.ReqCreateOrder, userID int64) error
// 	DeleteOrder(orderID int64, userID int64) error
// }

type OrderServices interface {
	CheckOut(req *request.ReqCheckout, userID int64) error
	CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) error
	GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error)
}
