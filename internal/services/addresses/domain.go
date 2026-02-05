package addressservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type AddressServices interface {
	GetUserAddressActive(userId int64) (*response.ResAddress, error)
	GetAllAddress(userid int64) ([]*response.ResAddress, error)
	CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error)
	UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error)
}
