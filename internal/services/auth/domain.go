package authservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type AuthServices interface {
	Login(req *request.ReqLogin) (*response.ResUser, error)
	Register(req *request.ReqCreateUser) (*response.ResUser, error)
}
