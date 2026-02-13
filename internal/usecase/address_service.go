package usecase

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

type AddressServiceImpl struct {
	domain.AddressRepositories
	*gorm.DB
}

func NewAddressService(address domain.AddressRepositories, db *gorm.DB) domain.AddressServices {
	return &AddressServiceImpl{
		AddressRepositories: address,
		DB:                  db,
	}
}

func (addressService *AddressServiceImpl) loadAddress(address *model.AddressModel) *response.ResAddress {
	return &response.ResAddress{
		ID:       address.ID,
		City:     address.City,
		Address:  address.Address,
		IsActive: address.IsActive,
	}
}
func (addressService *AddressServiceImpl) GetUserAddressActive(userId int64) (*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetUserAddressActive(addressService.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return addressService.loadAddress(result), nil
}
func (addressService *AddressServiceImpl) GetAllAddress(userId int64) ([]*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetAllAddress(addressService.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	addresss := []*response.ResAddress{}
	for _, address := range result {
		addresss = append(addresss, addressService.loadAddress(address))
	}

	return addresss, nil
}

func (addressService *AddressServiceImpl) CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error) {

	var addressModel *model.AddressModel
	err := addressService.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if _, err := addressService.AddressRepositories.GetUserAddressActive(tx, userID); err == nil {
				if err := addressService.AddressRepositories.DeActivate(tx, userID); err != nil {
					return pkg.MappingError(err)
				}
			}
		}

		result, err := addressService.AddressRepositories.CreateAddress(addressService.DB, &model.AddressModel{
			UserID:   userID,
			City:     request.City,
			Address:  request.Address,
			IsActive: *request.IsActive,
		})
		if err != nil {
			return pkg.MappingError(err)
		}

		addressModel = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	return addressService.loadAddress(addressModel), nil
}

func (addressService *AddressServiceImpl) UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error) {

	if err := addressService.AddressRepositories.CheckNotFoundForUpdate(addressService.DB, addressID); err != nil {
		return nil, pkg.MappingError(err)
	}
	var addressModel *model.AddressModel
	err := addressService.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if err := addressService.AddressRepositories.DeActivate(tx, userID); err != nil {
				return pkg.MappingError(err)
			}
		}

		result, err := addressService.AddressRepositories.UpdateAddressByUserId(tx, &model.AddressModel{
			ID:       addressID,
			UserID:   userID,
			City:     request.City,
			Address:  request.Address,
			IsActive: *request.IsActive,
		})

		if err != nil {
			return pkg.MappingError(err)
		}

		addressModel = result
		return nil
	})

	if err != nil {
		return nil, err
	}

	return addressService.loadAddress(addressModel), nil
}
