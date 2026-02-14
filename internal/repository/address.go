package repository

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"

	"gorm.io/gorm"
)

type AddressRepositoryImpl struct {
}

func NewAddressRepository() domain.AddressRepository {
	return &AddressRepositoryImpl{}
}

// PRIVATE
func (addressRepo *AddressRepositoryImpl) checkErrMysql(err error) error {
	if model.IsInternalErrMysql(err) {
		return err
	}
	if model.ForeignKeyErr(err) {
		return model.ErrUserNotFound
	}
	return model.ErrAddressNotFound
}
func (addressRepo *AddressRepositoryImpl) CheckNotFoundForUpdate(db *gorm.DB, addressID int64) error {
	var count int64
	if err := db.Model(&model.AddressModel{}).Where("id = ?", addressID).Count(&count).Error; err != nil {
		return model.ErrAddressNotFound
	}
	return nil
}

func (addressRepo *AddressRepositoryImpl) DeActivate(db *gorm.DB, userID int64) error {
	result := db.Model(&model.AddressModel{}).Where("user_id=?", userID).Update("is_active", false)
	if result.Error != nil {
		return addressRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil
	}

	return nil
}

// END PRIVATE

func (addressRepo *AddressRepositoryImpl) GetUserAddressActive(db *gorm.DB, userId int64) (*model.AddressModel, error) {
	var address model.AddressModel

	if err := db.Where("user_id=? AND is_active=?", userId, true).Take(&address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAddressById(db *gorm.DB, addressId int64) (*model.AddressModel, error) {
	var address model.AddressModel

	if err := db.Where("id=?", addressId).Take(&address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAllAddress(db *gorm.DB, userId int64) ([]*model.AddressModel, error) {
	var addresss []*model.AddressModel

	if err := db.Where("user_id=?", userId).Find(&addresss).Error; err != nil {
		return nil, err
	}

	return addresss, nil
}

func (addressRepo *AddressRepositoryImpl) CreateAddress(db *gorm.DB, address *model.AddressModel) (*model.AddressModel, error) {

	if err := db.Create(address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return addressRepo.GetAddressById(db, address.ID)
}

func (addressRepo *AddressRepositoryImpl) UpdateAddressByUserId(db *gorm.DB, address *model.AddressModel) (*model.AddressModel, error) {
	result := db.Model(&model.AddressModel{}).Where("id = ? AND user_id = ?", address.ID, address.UserID).Updates(address)
	if result.Error != nil {
		return nil, addressRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected < 1 {
		return nil, model.ErrAddressNotFound
	}
	return addressRepo.GetAddressById(db, address.ID)
}
