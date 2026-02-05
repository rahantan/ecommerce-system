package userservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
)

type UserServices interface {
	Login(req *request.ReqLogin) (*response.ResUser, error)
	Register(req *request.ReqCreateUser) (*response.ResUser, error)
}
