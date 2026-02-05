package orderrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type OrderRepositories interface {
	CreateOrder(db *gorm.DB, order *models.OrderModel) error
	GetAllOrder(db *gorm.DB, userID int64) ([]*models.OrderModel, error)
	GetAllOrderItem(db *gorm.DB, orderID, userID int64) ([]*models.OrderItemModel, error)
	DeleteOrder(db *gorm.DB, orderID int64, userID int64) error
	// CreateAddressOrder(db *gorm.DB, addressOrder *models.AddressOrderModel) error
}

type CheckOutRepositories interface {
	CheckOut(db *gorm.DB, checkout *models.CheckoutModel) error
	CheckOutConfirm(db *gorm.DB, checkout *models.CheckoutModel, userID int64) error
	GetLastDraftCheckOut(db *gorm.DB, userID int64) (*models.CheckoutModel, error)
	UpdateStatusLastCheckOut(db *gorm.DB, status string, userID int64) error
}
