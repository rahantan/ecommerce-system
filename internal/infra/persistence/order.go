package persistence

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
}

func NewOrderRepository() domain.OrderRepository {
	return &OrderRepositoryImpl{}
}
func (orderRepo *OrderRepositoryImpl) CreateOrder(db *gorm.DB, order *model.OrderModel) (*model.OrderModel, error) {
	if err := db.Create(order).Error; err != nil {
		return nil, err
	}
	return orderRepo.GetOrderByID(db, order.ID)
}

func (orderRepo *OrderRepositoryImpl) GetOrderByID(db *gorm.DB, orderID int64) (*model.OrderModel, error) {
	var order model.OrderModel
	if err := db.Where("id=? ", orderID).Preload("Payment").Preload("OrderItem").Preload("OrderItem.Product").Take(&order).Error; err != nil {
		return nil, model.ErrOrderNotFound
	}

	return &order, nil
}
func (orderRepo *OrderRepositoryImpl) GetAllOrderByUserID(db *gorm.DB, userID int64) ([]*model.OrderModel, error) {
	var orders []*model.OrderModel
	if err := db.Where("status_id != ?", 6).Preload("OrderItem").Preload("OrderStatus").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
func (orderRepo *OrderRepositoryImpl) GetAllOrder(db *gorm.DB) ([]*model.OrderModel, error) {
	var orders []*model.OrderModel
	if err := db.Preload("OrderItem").Preload("OrderStatus").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (orderRepo *OrderRepositoryImpl) GetOrderDetailsByID(db *gorm.DB, userID, orderID int64) (*model.OrderModel, error) {
	var order *model.OrderModel
	if err := db.Where("id=? AND user_id=?", orderID, userID).Preload("OrderStatus").Preload("OrderItem").Preload("OrderItem.Product").Take(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// func (orderRepo *OrderRepositoryImpl) GetAllOrderByStatusID(db *gorm.DB, userID, status int64) ([]*model.OrderModel, error) {
// 	var orders []*model.OrderModel
// 	if err := db.Where("user_id=?", userID).Preload("OrderStatus").Find(&orders).Error; err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

func (orderRepo *OrderRepositoryImpl) UpdateAll(db *gorm.DB, order *model.OrderModel) error {
	if err := db.Model(&model.OrderModel{}).Where("id=? ", order.ID).Updates(order).Error; err != nil {
		return err
	}
	return nil
}

func (orderRepo *OrderRepositoryImpl) UpdateStatusOrder(db *gorm.DB, orderID int64, statusID int64) error {
	if err := db.Model(&model.OrderModel{}).Where("id=? ", orderID).Update("status_id", statusID).Error; err != nil {
		return err
	}
	return nil
}

func (orderRepo *OrderRepositoryImpl) DeleteOrder(db *gorm.DB, orderID int64, userID int64) error {
	if err := db.Where("order_id=? AND user_id=?", orderID, userID).Delete(&model.OrderModel{}).Error; err != nil {
		return err
	}
	return nil
}
