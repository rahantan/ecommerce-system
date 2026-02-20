package user

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

func (userUC *UserUseCaseImpl) GetUserAddressActive(userId int64) (*response.ResAddress, error) {
	result, err := userUC.AddressRepository.GetUserAddressActive(userUC.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return userUC.loadAddress(result), nil
}

func (userUC *UserUseCaseImpl) GetAllAddress(userId int64) ([]*response.ResAddress, error) {
	result, err := userUC.AddressRepository.GetAllAddress(userUC.DB, userId)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	addresss := []*response.ResAddress{}
	for _, address := range result {
		addresss = append(addresss, userUC.loadAddress(address))
	}

	return addresss, nil
}

func (userUC *UserUseCaseImpl) CreateAddress(request *request.ReqCreateAddress, userID int64) (*response.ResAddress, error) {

	var addressModel *model.AddressModel
	err := userUC.DB.Transaction(func(tx *gorm.DB) error {

		if request.IsActive != nil && *request.IsActive {
			if _, err := userUC.AddressRepository.GetUserAddressActive(tx, userID); err == nil {
				if err := userUC.AddressRepository.DeActivate(tx, userID); err != nil {
					return pkg.MappingError(err)
				}
			}
		}

		result, err := userUC.AddressRepository.CreateAddress(tx, &model.AddressModel{
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

	return userUC.loadAddress(addressModel), nil
}

func (userUC *UserUseCaseImpl) UpdateAddressByUserId(request *request.ReqUpdateAddress, addressID int64, userID int64) (*response.ResAddress, error) {

	if err := userUC.AddressRepository.CheckNotFoundForUpdate(userUC.DB, addressID); err != nil {
		return nil, pkg.MappingError(err)
	}

	var addressModel *model.AddressModel
	err := userUC.DB.Transaction(func(tx *gorm.DB) error {

		if *request.IsActive {
			if err := userUC.AddressRepository.DeActivate(tx, userID); err != nil {
				return pkg.MappingError(err)
			}
		}

		result, err := userUC.AddressRepository.UpdateAddressByUserId(tx, &model.AddressModel{
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

	return userUC.loadAddress(addressModel), nil
}
