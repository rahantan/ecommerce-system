package orderrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type CheckOutRepositoryImpl struct {
}

func NewCheckoutSession() CheckOutRepositories {
	return &CheckOutRepositoryImpl{}
}
func (coRepo *CheckOutRepositoryImpl) CheckOut(db *gorm.DB, checkout *models.CheckoutModel) error {
	if err := db.Create(checkout).Error; err != nil {
		return err
	}
	return nil
}
func (coRepo *CheckOutRepositoryImpl) CheckOutConfirm(db *gorm.DB, checkout *models.CheckoutModel, userID int64) error {
	if err := db.Model(&models.CheckoutModel{}).Where("user_id=?", userID).Updates(checkout).Error; err != nil {
		return err
	}
	return nil
}

func (coRepo *CheckOutRepositoryImpl) GetLastDraftCheckOut(db *gorm.DB, userID int64) (*models.CheckoutModel, error) {
	var checkOutSession models.CheckoutModel

	result := db.Where("user_id=? AND status=?", userID, "draft").Preload("CheckoutItem").Last(&checkOutSession)
	if result.Error != nil {
		return nil, models.ErrCheckOutNotFound
	}

	return &checkOutSession, nil
}

func (coRepo *CheckOutRepositoryImpl) UpdateStatusLastCheckOut(db *gorm.DB, status string, userID int64) error {
	if err := db.Model(&models.CheckoutModel{}).Where("status=? AND user_id=? ", "draft", userID).Update("status", status).Error; err != nil {
		return nil
	}
	return nil
}
