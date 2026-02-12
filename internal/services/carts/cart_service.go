package cartitemsservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	cartitemrepositories "ecommerce-system/internal/repositories/carts"
	productservices "ecommerce-system/internal/services/products"
	"ecommerce-system/internal/utils"
	"fmt"

	"gorm.io/gorm"
)

type CartItemServiceImpl struct {
	*gorm.DB
	cartitemrepositories.CartItemRepositories
	productservices.ProductServices
}

func NewCartItemService(cartrepo cartitemrepositories.CartItemRepositories, productService productservices.ProductServices, db *gorm.DB) CartItemServices {
	return &CartItemServiceImpl{
		DB:                   db,
		CartItemRepositories: cartrepo,
		ProductServices:      productService,
	}
}

func (cartService *CartItemServiceImpl) loadRes(cartItem *models.CartItemModel) *response.ResCartItem {
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

func (cartService *CartItemServiceImpl) CreateOrUpdateCartItem(request *request.ReqCreateOrUpdateCartItem, userID int64) (*response.ResCartItem, error) {
	product, err := cartService.ProductServices.GetProductById(request.ProductID)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	result, err := cartService.CartItemRepositories.CreateOrUpdateCartItem(
		cartService.DB,
		&models.CartItemModel{
			UserID:    userID,
			ProductID: request.ProductID,
			Qty:       request.Qty,
			SubTotal:  (float64(request.Qty * int(product.Price))),
		},
	)

	if err != nil {
		return nil, utils.MappingError(err)
	}

	return cartService.loadRes(result), nil
}

func (cartService *CartItemServiceImpl) GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error) {
	results, err := cartService.CartItemRepositories.GetAllUserCartItem(cartService.DB, userID)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	items := []*response.ResCartItem{}
	for _, item := range results {
		items = append(items, cartService.loadRes(item))
	}

	return items, nil

}
func (cartService *CartItemServiceImpl) DeleteCartItemsByIDs(cartIDs []int64, userID int64) error {

	cartExist, err := cartService.CartItemRepositories.GetAllUserCartItem(cartService.DB, userID)
	if err != nil {
		return utils.MappingError(err)
	}

	cartIDexist := make(map[int64]bool)
	for _, cartItem := range cartExist {
		cartIDexist[cartItem.ID] = true
	}

	for _, cartID := range cartIDs {
		if !cartIDexist[cartID] {
			return exceptions.NewError(exceptions.KindNotFound, fmt.Sprintf("cart id %d not exsist", cartID), nil)
		}
	}

	if err := cartService.CartItemRepositories.DeleteCartItemsByIDs(cartService.DB, cartIDs, userID); err != nil {
		return utils.MappingError(err)
	}

	return nil
}

func (cartService *CartItemServiceImpl) GetCartItemByProductUser(productID int64, userID int64) (*response.ResCartItem, error) {

	result, err := cartService.CartItemRepositories.GetCartItemByProductUser(cartService.DB, productID, userID)
	if err != nil {
		return nil, utils.MappingError(err)
	}

	return cartService.loadRes(result), nil
}
