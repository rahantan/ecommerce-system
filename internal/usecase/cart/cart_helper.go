package cart

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
)

func (cartUC *CartItemUseCaseImpl) getUserCartMap(userID int64) (map[int64]*model.CartItemModel, error) {
	cartItem, err := cartUC.CartRepository.GetAllUserCartItem(cartUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	cartMap := make(map[int64]*model.CartItemModel)
	for _, item := range cartItem {
		cartMap[item.ProductID] = item
	}
	return cartMap, nil
}

func (cartUC *CartItemUseCaseImpl) resCart(cartItem *model.CartItemModel) *response.ResCartItem {
	return &response.ResCartItem{
		CartID: cartItem.ID,
		Product: response.ResProduct{
			ID:    cartItem.ProductID,
			Name:  cartItem.Product.Name,
			Price: cartItem.Product.Price,
			Stock: cartItem.Product.Stock,
		},
		Qty:      cartItem.Qty,
		SubTotal: cartItem.SubTotal,
	}
}
