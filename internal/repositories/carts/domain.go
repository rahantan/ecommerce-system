package cartitemrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type CartItemRepositories interface {
	CreateOrUpdateCartItem(db *gorm.DB, cartItem *models.CartItemModel) (*models.CartItemModel, error)
	GetAllUserCartItem(db *gorm.DB, userID int64) ([]*models.CartItemModel, error)
	GetAllCartItemByIDs(db *gorm.DB, cartIDs []int64, userID int64) ([]*models.CartItemModel, error)
	DeleteCartItemsByIDs(db *gorm.DB, cartIDs []int64, userID int64) error
	GetCartItemByProductUser(db *gorm.DB, cartID int64, userID int64) (*models.CartItemModel, error)
}
