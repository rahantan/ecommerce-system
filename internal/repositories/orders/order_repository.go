package orderrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
}

func NewOrderRepository() OrderRepositories {
	return &OrderRepositoryImpl{}
}
func (orderRepo *OrderRepositoryImpl) CreateOrder(db *gorm.DB, order *models.OrderModel) error {
	if err := db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (orderRepo *OrderRepositoryImpl) GetAllOrder(db *gorm.DB, userID int64) ([]*models.OrderModel, error) {
	var orders []*models.OrderModel
	if err := db.Where("userID=?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (orderRepo *OrderRepositoryImpl) GetAllOrderItem(db *gorm.DB, orderID, userID int64) ([]*models.OrderItemModel, error) {
	var orderItems []*models.OrderItemModel
	if err := db.Where("order_id=? AND userID=?", orderID, userID).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (orderRepo *OrderRepositoryImpl) DeleteOrder(db *gorm.DB, orderID int64, userID int64) error {
	if err := db.Where("order_id=? AND user_id=?", orderID, userID).Delete(&models.OrderModel{}).Error; err != nil {
		return err
	}
	return nil
}

// func (orderRepo *OrderRepositoryImpl) CreateAddressOrder(db *gorm.DB, addressOrder *models.AddressOrderModel) error {
// 	if err := db.Create(addressOrder).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
