package addressservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	addressrepositories "ecommerce-system/internal/repositories/addresses"

	"gorm.io/gorm"
)

type AddressServiceImpl struct {
	addressrepositories.AddressRepositories
	*gorm.DB
}

func NewAddressService(address addressrepositories.AddressRepositories, db *gorm.DB) AddressServices {
	return &AddressServiceImpl{
		AddressRepositories: address,
		DB:                  db,
	}
}

func (addressService *AddressServiceImpl) handleError(err error) error {
	return exceptions.CheckError(err)
}

func (addressService *AddressServiceImpl) loadAddress(address *models.AddressModel) *response.ResAddress {
	return &response.ResAddress{
		ID:       address.ID,
		City:     address.City,
		Address:  address.Address,
		IsActive: address.IsActive,
	}
}
func (addressService *AddressServiceImpl) GetUserAddressActive(userId int64) (*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetUserAddressActive(addressService.DB, userId)
	if errCheck := addressService.handleError(err); errCheck != nil {
		return nil, errCheck
	}
	return addressService.loadAddress(result), nil
}
func (addressService *AddressServiceImpl) GetAllAddress(userId int64) ([]*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetAllAddress(addressService.DB, userId)
	if errCheck := addressService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	addresss := []*response.ResAddress{}
	for _, address := range result {
		addresss = append(addresss, addressService.loadAddress(address))
	}

	return addresss, nil
}

func (addressService *AddressServiceImpl) CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error) {

	var addressModel *models.AddressModel
	err := addressService.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if _, err := addressService.AddressRepositories.GetUserAddressActive(tx, userID); err == nil {
				if err := addressService.AddressRepositories.DeActivate(tx, userID); err != nil {
					return err
				}
			}
		}

		result, err := addressService.AddressRepositories.CreateAddress(addressService.DB, &models.AddressModel{
			UserID:   userID,
			City:     request.City,
			Address:  request.Address,
			IsActive: *request.IsActive,
		})
		if err != nil {
			return err
		}

		addressModel = result
		return nil
	})

	if err != nil {
		return nil, addressService.handleError(err)
	}
	return addressService.loadAddress(addressModel), nil
}

func (addressService *AddressServiceImpl) UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error) {

	if err := addressService.AddressRepositories.CheckNotFoundForUpdate(addressService.DB, addressID); err != nil {
		return nil, addressService.handleError(err)
	}
	var addressModel *models.AddressModel
	err := addressService.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if err := addressService.AddressRepositories.DeActivate(tx, userID); err != nil {
				return err
			}
		}

		result, err := addressService.AddressRepositories.UpdateAddressByUserId(tx, &models.AddressModel{
			ID:       addressID,
			UserID:   userID,
			City:     request.City,
			Address:  request.Address,
			IsActive: *request.IsActive,
		})

		if err != nil {
			return err
		}

		addressModel = result
		return nil
	})

	if err != nil {
		return nil, addressService.handleError(err)
	}

	return addressService.loadAddress(addressModel), nil
}
