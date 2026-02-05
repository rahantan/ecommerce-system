package addressrepositories

import (
	"ecommerce-system/internal/models"

	"gorm.io/gorm"
)

type AddressRepositories interface {
	GetAddressById(db *gorm.DB, addressID int64) (*models.AddressModel, error)
	GetUserAddressActive(db *gorm.DB, userId int64) (*models.AddressModel, error)
	GetAllAddress(db *gorm.DB, userId int64) ([]*models.AddressModel, error)
	UpdateAddressByUserId(db *gorm.DB, address *models.AddressModel) (*models.AddressModel, error)
	CreateAddress(db *gorm.DB, address *models.AddressModel) (*models.AddressModel, error)
	CheckNotFoundForUpdate(db *gorm.DB, addressID int64) error
	DeActivate(db *gorm.DB, userID int64) error
	// ActivateAddress(addressID int64, userID int64) error
}
