package addressrepositories

import (
	"ecommerce-system/internal/models"

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

// PRIVATE
func (addressRepo *AddressRepositoryImpl) checkErrMysql(err error) error {
	if models.IsInternalErrMysql(err) {
		return err
	}
	if models.ForeignKeyErr(err) {
		return models.ErrUserNotFound
	}
	return models.ErrAddressNotFound
}
func (addressRepo *AddressRepositoryImpl) checkNotFoundForUpdate(addressID int64) bool {
	var count int64
	if err := addressRepo.DB.Model(&models.AddressModel{}).Where("id = ?", addressID).Count(&count).Error; err != nil {
		return false
	}
	return count == 0
}

func (addressRepo *AddressRepositoryImpl) deActivate(userID int64) error {
	return addressRepo.DB.Model(&models.AddressModel{UserID: userID}).Update("is_activate=", false).Error
}

// END PRIVATE

func (addressRepo *AddressRepositoryImpl) ActivateAddress(addressID int64, userID int64) error {
	if err := addressRepo.deActivate(userID); err != nil {
		return addressRepo.checkErrMysql(err)
	}

	if addressRepo.checkNotFoundForUpdate(addressID) {
		return models.ErrAddressNotFound
	}

	if err := addressRepo.DB.Model(&models.AddressModel{ID: addressID, UserID: userID}).Update("is_activate=", true).Error; err != nil {
		return err
	}

	return nil
}

func (addressRepo *AddressRepositoryImpl) GetUserActiveAddress(userId int64) (*models.AddressModel, error) {
	var address models.AddressModel

	if err := addressRepo.DB.Where("user_id=? AND is_active=?", userId, true).Take(&address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAddressById(addressId int64) (*models.AddressModel, error) {
	var address models.AddressModel

	err := addressRepo.DB.Where("id=?", addressId).Take(&address).Error
	if err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAllAddress(userId int64) ([]*models.AddressModel, error) {
	var addresss []*models.AddressModel

	if err := addressRepo.DB.Where("user_id=?", userId).Find(&addresss).Error; err != nil {
		return nil, err
	}

	return addresss, nil
}
func (addressRepo *AddressRepositoryImpl) UpdateAddressByUserId(address *models.AddressModel) (*models.AddressModel, error) {

	if addressRepo.checkNotFoundForUpdate(address.ID) {
		return nil, models.ErrAddressNotFound
	}

	result := addressRepo.DB.Model(&models.AddressModel{}).Where("id = ? AND user_id = ?", address.ID, address.UserID).Updates(address)
	if result.Error != nil {
		return nil, addressRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected < 1 {
		return nil, models.ErrNoRowsAffected
	}

	return addressRepo.GetAddressById(address.ID)
}

func (addressRepo *AddressRepositoryImpl) CreateAddress(address *models.AddressModel) (*models.AddressModel, error) {

	if err := addressRepo.DB.Create(address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return addressRepo.GetAddressById(address.ID)
}
