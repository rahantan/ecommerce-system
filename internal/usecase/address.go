package usecase

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

type AddressUseCaseImpl struct {
	domain.AddressRepository
	*gorm.DB
}

func NewAddressUseCase(address domain.AddressRepository, db *gorm.DB) domain.AddressUseCase {
	return &AddressUseCaseImpl{
		AddressRepository: address,
		DB:                db,
	}
}

func (addressUC *AddressUseCaseImpl) loadAddress(address *model.AddressModel) *response.ResAddress {
	return &response.ResAddress{
		ID:       address.ID,
		City:     address.City,
		Address:  address.Address,
		IsActive: address.IsActive,
	}
}
func (addressUC *AddressUseCaseImpl) GetUserAddressActive(userId int64) (*response.ResAddress, error) {
	result, err := addressUC.AddressRepository.GetUserAddressActive(addressUC.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return addressUC.loadAddress(result), nil
}
func (addressUC *AddressUseCaseImpl) GetAllAddress(userId int64) ([]*response.ResAddress, error) {
	result, err := addressUC.AddressRepository.GetAllAddress(addressUC.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	addresss := []*response.ResAddress{}
	for _, address := range result {
		addresss = append(addresss, addressUC.loadAddress(address))
	}

	return addresss, nil
}

func (addressUC *AddressUseCaseImpl) CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error) {

	var addressModel *model.AddressModel
	err := addressUC.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if _, err := addressUC.AddressRepository.GetUserAddressActive(tx, userID); err == nil {
				if err := addressUC.AddressRepository.DeActivate(tx, userID); err != nil {
					return pkg.MappingError(err)
				}
			}
		}

		result, err := addressUC.AddressRepository.CreateAddress(addressUC.DB, &model.AddressModel{
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

	return addressUC.loadAddress(addressModel), nil
}

func (addressUC *AddressUseCaseImpl) UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error) {

	if err := addressUC.AddressRepository.CheckNotFoundForUpdate(addressUC.DB, addressID); err != nil {
		return nil, pkg.MappingError(err)
	}
	var addressModel *model.AddressModel
	err := addressUC.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if err := addressUC.AddressRepository.DeActivate(tx, userID); err != nil {
				return pkg.MappingError(err)
			}
		}

		result, err := addressUC.AddressRepository.UpdateAddressByUserId(tx, &model.AddressModel{
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

	return addressUC.loadAddress(addressModel), nil
}
