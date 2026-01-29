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

func (addressRepo *AddressRepositoryImpl) deActivate(tx *gorm.DB, userID int64) error {
	result := tx.Model(&models.AddressModel{}).Where("user_id=?", userID).Update("is_active", false)
	if result.Error != nil {
		return addressRepo.checkErrMysql(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil
	}

	return nil
}

// END PRIVATE

func (addressRepo *AddressRepositoryImpl) GetUserActiveAddress(userId int64) (*models.AddressModel, error) {
	var address models.AddressModel

	if err := addressRepo.DB.Where("user_id=? AND is_active=?", userId, true).Take(&address).Error; err != nil {
		return nil, addressRepo.checkErrMysql(err)
	}

	return &address, nil
}
func (addressRepo *AddressRepositoryImpl) GetAddressById(addressId int64) (*models.AddressModel, error) {
	var address models.AddressModel

	if err := addressRepo.DB.Where("id=?", addressId).Take(&address).Error; err != nil {
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
	// Check if address exists first
	if addressRepo.checkNotFoundForUpdate(address.ID) {
		return nil, models.ErrAddressNotFound
	}

	// Use transaction for consistent state
	err := addressRepo.Transaction(func(tx *gorm.DB) error {

		if address.IsActive {
			if err := addressRepo.deActivate(tx, address.ID); err != nil {
				return addressRepo.checkErrMysql(err)
			}
		}

		result := tx.Model(&models.AddressModel{}).Where("id = ? AND user_id = ?", address.ID, address.UserID).Updates(address)
		if result.Error != nil {
			return addressRepo.checkErrMysql(result.Error)
		}

		if result.RowsAffected < 1 {
			return models.ErrNoRowsAffected
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return addressRepo.GetAddressById(address.ID)
}

func (addressRepo *AddressRepositoryImpl) CreateAddress(address *models.AddressModel) (*models.AddressModel, error) {

	err := addressRepo.Transaction(func(tx *gorm.DB) error {
		if address.IsActive {
			if err := addressRepo.deActivate(tx, address.ID); err != nil {
				return addressRepo.checkErrMysql(err)
			}
		}
		if err := tx.Create(address).Error; err != nil {
			return addressRepo.checkErrMysql(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return addressRepo.GetAddressById(address.ID)
}

// func (addressRepo *AddressRepositoryImpl) ActivateAddress(addressID int64, userID int64) error {
// 	if err := addressRepo.deActivate(userID); err != nil {
// 		return addressRepo.checkErrMysql(err)
// 	}

// 	if addressRepo.checkNotFoundForUpdate(addressID) {
// 		return models.ErrAddressNotFound
// 	}

// 	if err := addressRepo.DB.Model(&models.AddressModel{ID: addressID, UserID: userID}).Update("is_activate=", true).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }
