package addressservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	addressrepositories "ecommerce-system/internal/repositories/addresses"
)

type AddressServiceImpl struct {
	addressrepositories.AddressRepositories
}

func NewAddresyService(address addressrepositories.AddressRepositories) AddressServices {
	return &AddressServiceImpl{
		AddressRepositories: address,
	}
}

func (addressService *AddressServiceImpl) loadAddress(address *models.AddressModel) *response.ResAddress {
	return &response.ResAddress{
		ID:      address.ID,
		City:    address.City,
		Address: address.Address,
	}
}
func (addressService *AddressServiceImpl) GetUserActiveAddress(userId int64) (*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetUserActiveAddress(userId)
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return addressService.loadAddress(result), nil
}
func (addressService *AddressServiceImpl) GetAllAddress(userId int64) ([]*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.GetAllAddress(userId)
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}

	addresss := []*response.ResAddress{}
	for _, address := range result {
		addresss = append(addresss, addressService.loadAddress(address))
	}

	return addresss, nil
}

func (addressService *AddressServiceImpl) CreateAddress(request *request.ReqCreateAddress) (*response.ResAddress, error) {
	result, err := addressService.AddressRepositories.CreateAddress(&models.AddressModel{
		UserID:   request.UserID,
		City:     request.City,
		Address:  request.Address,
		IsActive: request.IsActive,
	})
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return addressService.loadAddress(result), nil
}
func (addressService *AddressServiceImpl) UpdateAddressByUserId(request *request.ReqUpdateAddress) (*response.ResAddress, error) {

	result, err := addressService.AddressRepositories.UpdateAddressByUserId(&models.AddressModel{
		ID:       request.ID,
		UserID:   request.UserID,
		City:     request.City,
		Address:  request.Address,
		IsActive: request.IsActive,
	})
	if errCheck := exceptions.CheckError(err); errCheck != nil {
		return nil, errCheck
	}
	return addressService.loadAddress(result), nil
}
