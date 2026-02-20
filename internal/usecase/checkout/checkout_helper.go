package checkout

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"fmt"

	"gorm.io/gorm"
)

func (coUC *CheckOutUseCaseImpl) cartItems(cartIDs []int64, userID int64) (int64, []model.CheckoutItemModel, error) {

	items, err := coUC.CartRepository.GetAllCartItemByIDs(coUC.DB, cartIDs, userID)
	if err != nil {
		return 0, nil, pkg.MappingError(err)
	}

	if len(cartIDs) != len(items) {
		return 0, nil, pkg.NewError(pkg.KindNotFound, "cart item or product not found", nil)
	}

	var (
		checkOutItems []model.CheckoutItemModel
		totalPrice    int64
	)

	for _, item := range items {

		checkOutItems = append(checkOutItems, model.CheckoutItemModel{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     item.Price,
			SubTotal:  item.SubTotal,
			CartID:    &item.ID,
		})

		totalPrice += item.SubTotal
	}

	return totalPrice, checkOutItems, nil
}

func (coUC *CheckOutUseCaseImpl) directItems(reqItems []request.ReqItem) (int64, []model.CheckoutItemModel, error) {

	var (
		checkOutItems []model.CheckoutItemModel
		productIDs    []int64
		totalPrice    int64
	)

	for _, item := range reqItems {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := coUC.ProductRepository.FindByIDs(coUC.DB, productIDs)
	if err != nil {
		return 0, nil, err
	}

	for _, item := range reqItems {

		subTotal := products[item.ProductID].Price * int64(item.Qty)
		checkOutItems = append(checkOutItems, model.CheckoutItemModel{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     products[item.ProductID].Price,
			SubTotal:  subTotal,
		})

		totalPrice += subTotal

	}

	return totalPrice, checkOutItems, nil
}

func (coUC *CheckOutUseCaseImpl) resolveAddress(tx *gorm.DB, addressID int64, userID int64) (*model.AddressModel, error) {

	if addressID == 0 {
		return coUC.AddressRepository.GetUserAddressActive(tx, userID)
	}

	return coUC.AddressRepository.GetAddressById(tx, addressID)
}

func (coUC *CheckOutUseCaseImpl) deleteCartItems(tx *gorm.DB, items []model.CheckoutItemModel, userID int64) error {

	var cartIDs []int64
	for _, item := range items {
		if item.CartID != nil {
			cartIDs = append(cartIDs, *item.CartID)
		}
	}

	err := coUC.CartRepository.DeleteCartItemsByIDs(tx, cartIDs, userID)
	if err != nil {
		return pkg.MappingError(err)
	}

	return nil

}

func (coUC *CheckOutUseCaseImpl) validateItems(products map[int64]*model.ProductModel, items []model.CheckoutItemModel) error {

	for i := range items {

		product, ok := products[items[i].ProductID]
		if !ok {
			return pkg.NewError(pkg.KindNotFound, fmt.Sprintf("product id %d not found", items[i].ProductID), nil)
		}

		if items[i].Qty > product.Stock {
			return pkg.NewError(pkg.KindConflict, "qty cannot exceed available stock", nil)
		}

		items[i].Price = product.Price
		items[i].SubTotal = product.Price * int64(items[i].Qty)

	}

	return nil
}
