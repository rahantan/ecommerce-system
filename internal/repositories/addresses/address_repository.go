package addressrepositories

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	"errors"

	"gorm.io/gorm"
)

type AddressRepositoryImpl struct {
	*gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepositories {
	return &AddressRepositoryImpl{
		DB: db,
	}
}
func (addressRepo *AddressRepositoryImpl) GetUserActiveAddress(userId int64) (*models.AddressModel, error) {
	var address models.AddressModel
	err := addressRepo.DB.Where("user_id=?", userId).Where("is_active=?", true).Take(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrAddressNotFound
		}
		return nil, err
	}
	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAddressById(addressId int64) (*models.AddressModel, error) {
	var address models.AddressModel
	err := addressRepo.DB.Where("id=?", addressId).Take(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrAddressNotFound
		}
		return nil, err
	}
	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAllAddress(userId int64) ([]*models.AddressModel, error) {
	var addresss []*models.AddressModel
	err := addressRepo.DB.Where("user_id=?", userId).Find(&addresss).Error
	if err != nil {
		return nil, err
	}
	return addresss, nil
}
func (addressRepo *AddressRepositoryImpl) UpdateAddressByUserId(address *models.AddressModel) (*models.AddressModel, error) {
	result := addressRepo.DB.Model(&models.AddressModel{}).Where("id=?", address.ID).Where("user_id=?", address.UserID).Updates(address)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 1 {
		return nil, exceptions.ErrNoRowsAffected
	}
	return addressRepo.GetAddressById(address.ID)
}
func (addressRepo *AddressRepositoryImpl) CreateAddress(address *models.AddressModel) (*models.AddressModel, error) {
	err := addressRepo.DB.Create(address).Error
	if err != nil {
		return nil, err
	}
	return addressRepo.GetAddressById(address.ID)
}
