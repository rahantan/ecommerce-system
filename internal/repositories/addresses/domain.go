package addressrepositories

import "ecommerce-system/internal/models"

type AddressRepositories interface {
	GetAddressById(userId int64) (*models.AddressModel, error)
	GetUserActiveAddress(userId int64) (*models.AddressModel, error)
	GetAllAddress(userId int64) ([]*models.AddressModel, error)
	UpdateAddressByUserId(address *models.AddressModel) (*models.AddressModel, error)
	CreateAddress(address *models.AddressModel) (*models.AddressModel, error)
	// ActivateAddress(addressID int64, userID int64) error
}
