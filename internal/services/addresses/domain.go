package addressservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type AddressServices interface {
	UpdateAddressByUserId(request *request.ReqUpdateAddress) (*response.ResAddress, error)
	GetUserActiveAddress(userId int64) (*response.ResAddress, error)
	GetAllAddress(userid int64) ([]*response.ResAddress, error)
	CreateAddress(request *request.ReqCreateAddress) (*response.ResAddress, error)
	// ActivateAddress(addressID int64, userID int64) error
}
