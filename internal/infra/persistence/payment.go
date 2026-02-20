package persistence

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepositoryImpl struct {
}

func NewPaymentRepository() domain.PaymentRepository {
	return &PaymentRepositoryImpl{}
}

func (paymentRepo *PaymentRepositoryImpl) GetPaymentByOrderID(db *gorm.DB, orderID int64) (*model.PaymentOrderModel, error) {
	var payment model.PaymentOrderModel
	if err := db.Where("order_id=?", orderID).Last(&payment).Error; err != nil {
		return nil, model.ErrPaymentNotFound
	}
	return &payment, nil
}
func (paymentRepo *PaymentRepositoryImpl) UpdateStatusPayment(db *gorm.DB, orderID int64, status string) error {
	if err := db.Model(&model.PaymentOrderModel{}).Where("order_id=?", orderID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
func (paymentRepo *PaymentRepositoryImpl) SavePayment(db *gorm.DB, payment *model.PaymentOrderModel) error {
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "order_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"snap_token", "redirect_url", "updated_at", "status"}),
	}).Create(&payment).Error; err != nil {

		return err
	}
	return nil
}
