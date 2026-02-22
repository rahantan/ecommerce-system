package persistence

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartItemRepositoryImpl struct {
}

func NewCartRepository() domain.CartRepository {
	return &CartItemRepositoryImpl{}
}

func (cartRepo *CartItemRepositoryImpl) checkErrMysql(err error) error {
	if model.IsInternalErrMysql(err) {
		return err
	}
	if model.IsDuplicateKeyError(err) {
		return err
	}
	if model.ForeignKeyErr(err) {
		return model.ErrProductNotFound
	}

	return model.ErrCartItemNotFound
}

func (cartRepo *CartItemRepositoryImpl) checkNotFoundForUpdate(db *gorm.DB, cartID int64) bool {
	var count int64
	if err := db.Model(&model.CartItemModel{}).Where("id = ?", cartID).Count(&count).Error; err != nil {
		return false
	}
	return count == 0
}

func (cartRepo *CartItemRepositoryImpl) GetUserCartByID(db *gorm.DB, cartID, userID int64) (*model.CartItemModel, error) {
	var cartItem model.CartItemModel

	if err := db.Where("id=? AND user_id=?", cartID, userID).Preload("Product").Take(&cartItem).Error; err != nil {
		return nil, cartRepo.checkErrMysql(err)
	}
	return &cartItem, nil
}

func (cartRepo *CartItemRepositoryImpl) GetAllUserCartItem(db *gorm.DB, userID int64) ([]*model.CartItemModel, error) {
	var cartItems []*model.CartItemModel

	if err := db.Where("user_id=?", userID).Preload("Product").Find(&cartItems).Error; err != nil {
		return nil, err
	}

	return cartItems, nil
}

func (cartRepo *CartItemRepositoryImpl) GetAllCartItemByIDs(db *gorm.DB, cartIDs []int64, userID int64) ([]*model.CartItemModel, error) {
	var cartItems []*model.CartItemModel

	if err := db.Where("id IN ? AND user_id=?", cartIDs, userID).Preload("Product").Find(&cartItems).Error; err != nil {
		return nil, err
	}

	return cartItems, nil
}
func (cartRepo *CartItemRepositoryImpl) CreateOrUpdateCartItem(db *gorm.DB, cartItem *model.CartItemModel) (*model.CartItemModel, error) {

	if err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cartItem).Error; err != nil {
		return nil, cartRepo.checkErrMysql(err)
	}

	return cartRepo.GetUserCartByID(db, cartItem.ID, cartItem.UserID)
}

func (cartRepo *CartItemRepositoryImpl) DeleteCartItemsByIDs(db *gorm.DB, cartIDs []int64, userID int64) error {

	if err := db.Where("id IN ? AND user_id=?", cartIDs, userID).Delete(&model.CartItemModel{}).Error; err != nil {
		return cartRepo.checkErrMysql(err)
	}

	return nil
}
