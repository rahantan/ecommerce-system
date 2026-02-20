package cart

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"fmt"

	"gorm.io/gorm"
)

type CartItemUseCaseImpl struct {
	*gorm.DB
	domain.CartRepository
	domain.ProductUseCase
}

func NewCartUseCase(cartrepo domain.CartRepository, productService domain.ProductUseCase, db *gorm.DB) domain.CartUseCase {
	return &CartItemUseCaseImpl{
		DB:             db,
		CartRepository: cartrepo,
		ProductUseCase: productService,
	}
}

func (cartUC *CartItemUseCaseImpl) GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error) {
	results, err := cartUC.CartRepository.GetAllUserCartItem(cartUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	items := []*response.ResCartItem{}
	for _, item := range results {
		items = append(items, cartUC.resCart(item))
	}

	return items, nil

}

func (cartUC *CartItemUseCaseImpl) UpdateCartItemByID(request *request.ReqUpdateCartQty, cartID, userID int64) (*response.ResCartItem, error) {
	userCart, err := cartUC.CartRepository.GetUserCartByID(cartUC.DB, cartID, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	cart := &model.CartItemModel{
		ID:        userCart.ID,
		UserID:    userID,
		Qty:       request.Qty,
		ProductID: userCart.ProductID,
		Price:     userCart.Product.Price,
		SubTotal:  userCart.Product.Price * int64(request.Qty),
	}

	result, err := cartUC.CartRepository.CreateOrUpdateCartItem(cartUC.DB, cart)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return cartUC.resCart(result), nil
}

func (cartUC *CartItemUseCaseImpl) AddCartItem(request *request.ReqAddCart, userID int64) (*response.ResCartItem, error) {
	product, err := cartUC.ProductUseCase.GetProductById(request.ProductID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	userCartMap, err := cartUC.getUserCartMap(userID)
	if err != nil {
		return nil, err
	}

	cart := &model.CartItemModel{
		UserID:    userID,
		ProductID: request.ProductID,
		Price:     product.Price,
	}

	if len(userCartMap) < 1 {
		cart.Qty = 1
		cart.SubTotal = 1 * product.Price
	} else {
		qty := userCartMap[product.ID].Qty
		cart.Qty = qty + 1
		cart.SubTotal = int64(cart.Qty) * product.Price
	}

	result, err := cartUC.CartRepository.CreateOrUpdateCartItem(cartUC.DB, cart)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return cartUC.resCart(result), nil
}

func (cartUC *CartItemUseCaseImpl) DeleteCartItemsByIDs(cartIDs []int64, userID int64) error {

	cartExist, err := cartUC.CartRepository.GetAllUserCartItem(cartUC.DB, userID)
	if err != nil {
		return pkg.MappingError(err)
	}

	cartIDexist := make(map[int64]bool)
	for _, cartItem := range cartExist {
		cartIDexist[cartItem.ID] = true
	}

	for _, cartID := range cartIDs {
		if !cartIDexist[cartID] {
			return pkg.NewError(pkg.KindNotFound, fmt.Sprintf("cart id %d not exsist", cartID), nil)
		}
	}

	if err := cartUC.CartRepository.DeleteCartItemsByIDs(cartUC.DB, cartIDs, userID); err != nil {
		return pkg.MappingError(err)
	}

	return nil
}
