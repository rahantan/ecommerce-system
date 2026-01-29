package cartitemrepositories

import (
	"ecommerce-system/internal/models"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartItemRepositoryImpl struct {
}

func NewCartItemRepository() CartItemRepositories {
	return &CartItemRepositoryImpl{}
}

func (cartRepo *CartItemRepositoryImpl) GetCartItemByProductUser(db *gorm.DB, productID int64, userID int64) (*models.CartItemModel, error) {
	var cartItem models.CartItemModel

	if err := db.Where("product_id=? AND user_id=?", productID, userID).Preload("Product").Take(&cartItem).Error; err != nil {
		return nil, cartRepo.checkErrMysql(err)
	}
	return &cartItem, nil
}
func (cartRepo *CartItemRepositoryImpl) GetAllUserCartItem(db *gorm.DB, userID int64) ([]*models.CartItemModel, error) {
	var cartItems []*models.CartItemModel

	if err := db.Where("user_id=?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		return nil, err
	}

	return cartItems, nil
}
func (cartRepo *CartItemRepositoryImpl) CreateOrUpdateCartItem(db *gorm.DB, cartItem *models.CartItemModel) (*models.CartItemModel, error) {

	if err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cartItem).Error; err != nil {
		return nil, cartRepo.checkErrMysql(err)
	}
	fmt.Println("done")
	return cartRepo.GetCartItemByProductUser(db, cartItem.ProductID, cartItem.UserID)
}

func (cartRepo *CartItemRepositoryImpl) DeleteCartItemById(db *gorm.DB, cartID int64, userID int64) error {
	if err := db.Where("id=? AND user_id=?", cartID, userID).Delete(&models.CartItemModel{}).Error; err != nil {
		return cartRepo.checkErrMysql(err)
	}
	return nil
}

func (cartRepo *CartItemRepositoryImpl) checkErrMysql(err error) error {
	if models.IsInternalErrMysql(err) {
		return err
	}
	if models.IsDuplicateKeyError(err) {
		return err
	}
	if models.ForeignKeyErr(err) {
		return models.ErrProductNotFound
	}

	return models.ErrCartItemNotFound
}

func (cartRepo *CartItemRepositoryImpl) checkNotFoundForUpdate(db *gorm.DB, cartID int64) bool {
	var count int64
	if err := db.Model(&models.CartItemModel{}).Where("id = ?", cartID).Count(&count).Error; err != nil {
		return false
	}
	return count == 0
}
