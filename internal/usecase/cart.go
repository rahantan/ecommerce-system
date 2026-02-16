package usecase

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

func (cartUC *CartItemUseCaseImpl) loadRes(cartItem *model.CartItemModel) *response.ResCartItem {
	return &response.ResCartItem{
		CartItemID: cartItem.ID,
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

func (cartUC *CartItemUseCaseImpl) CreateOrUpdateCartItem(request *request.ReqCreateOrUpdateCartItem, userID int64) (*response.ResCartItem, error) {
	product, err := cartUC.ProductUseCase.GetProductById(request.ProductID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	result, err := cartUC.CartRepository.CreateOrUpdateCartItem(
		cartUC.DB,
		&model.CartItemModel{
			UserID:    userID,
			ProductID: request.ProductID,
			Qty:       request.Qty,
			SubTotal:  int64(request.Qty) * product.Price,
			Price:     product.Price,
		},
	)

	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return cartUC.loadRes(result), nil
}

func (cartUC *CartItemUseCaseImpl) GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error) {
	results, err := cartUC.CartRepository.GetAllUserCartItem(cartUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	items := []*response.ResCartItem{}
	for _, item := range results {
		items = append(items, cartUC.loadRes(item))
	}

	return items, nil

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

func (cartUC *CartItemUseCaseImpl) GetCartItemByProductUser(productID int64, userID int64) (*response.ResCartItem, error) {

	result, err := cartUC.CartRepository.GetCartItemByProductUser(cartUC.DB, productID, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	return cartUC.loadRes(result), nil
}
