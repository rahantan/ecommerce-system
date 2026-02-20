package user

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
)

func (userUC *UserUseCaseImpl) loadAddress(address *model.AddressModel) *response.ResAddress {
	return &response.ResAddress{
		ID:       address.ID,
		City:     address.City,
		Address:  address.Address,
		IsActive: address.IsActive,
	}
}
func (userUC *UserUseCaseImpl) loadUserRes(userMdl *model.UserModel) *response.ResUser {
	// addresses := []response.ResAddress{}
	// for _, address := range userMdl.Address {
	// 	addresses = append(addresses, response.ResAddress{
	// 		ID:      address.ID,
	// 		City:    address.City,
	// 		Address: address.Address,
	// 	})
	// }

	return &response.ResUser{
		ID:    userMdl.ID,
		Name:  userMdl.Name,
		Email: userMdl.Email,
		Phone: userMdl.Phone,
		Role: response.ResRole{
			ID:    userMdl.Role.ID,
			Title: userMdl.Role.Title,
		},
		// Addresses: addresses,
		CreatedAt: userMdl.CreatedAt.Format(pkg.DateTimeLayout),
		UpdatedAt: userMdl.UpdatedAt.Format(pkg.DateTimeLayout),
	}
}
