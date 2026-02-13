package repository

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
}

func NewOrderRepository() domain.OrderRepositories {
	return &OrderRepositoryImpl{}
}
func (orderRepo *OrderRepositoryImpl) CreateOrder(db *gorm.DB, order *model.OrderModel) error {
	if err := db.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func (orderRepo *OrderRepositoryImpl) GetAllOrder(db *gorm.DB, userID int64) ([]*model.OrderModel, error) {
	var orders []*model.OrderModel
	if err := db.Where("userID=?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (orderRepo *OrderRepositoryImpl) GetAllOrderItem(db *gorm.DB, orderID, userID int64) ([]*model.OrderItemModel, error) {
	var orderItems []*model.OrderItemModel
	if err := db.Where("order_id=? AND userID=?", orderID, userID).Find(&orderItems).Error; err != nil {
		return nil, err
	}
	return orderItems, nil
}

func (orderRepo *OrderRepositoryImpl) DeleteOrder(db *gorm.DB, orderID int64, userID int64) error {
	if err := db.Where("order_id=? AND user_id=?", orderID, userID).Delete(&model.OrderModel{}).Error; err != nil {
		return err
	}
	return nil
}

// func (orderRepo *OrderRepositoryImpl) CreateAddressOrder(db *gorm.DB, addressOrder *model.AddressOrderModel) error {
// 	if err := db.Create(addressOrder).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }
