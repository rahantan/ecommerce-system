package domain

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetUserCartByID(db *gorm.DB, cartID, userID int64) (*model.CartItemModel, error)
	CreateOrUpdateCartItem(db *gorm.DB, cartItem *model.CartItemModel) (*model.CartItemModel, error)
	GetAllUserCartItem(db *gorm.DB, userID int64) ([]*model.CartItemModel, error)
	GetAllCartItemByIDs(db *gorm.DB, cartIDs []int64, userID int64) ([]*model.CartItemModel, error)
	DeleteCartItemsByIDs(db *gorm.DB, cartIDs []int64, userID int64) error
	// GetCartItemByProductUser(db *gorm.DB, cartID int64, userID int64) (*model.CartItemModel, error)
}

type CartUseCase interface {
	// CreateOrUpdateCartItem(request *request.ReqCreateOrUpdateCartItem, userID int64) (*response.ResCartItem, error)
	GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error)
	DeleteCartItemsByIDs(cartIDs []int64, userID int64) error
	// GetCartItemByProductUser(productID int64, userID int64) (*response.ResCartItem, error)
	AddCartItem(request *request.ReqAddCart, userID int64) (*response.ResCartItem, error)
	UpdateCartItemByID(request *request.ReqUpdateCartQty, cartID, userID int64) (*response.ResCartItem, error)
}

type CartHandler interface {
	AddCartItem(ctx *fiber.Ctx) error
	GetAllUserCartItem(ctx *fiber.Ctx) error
	DeleteCartItemsByIDs(ctx *fiber.Ctx) error
	UpdateCartItemByID(ctx *fiber.Ctx) error
}
