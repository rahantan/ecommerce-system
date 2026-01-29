package cartitemsservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	cartitemrepositories "ecommerce-system/internal/repositories/carts"
	productservices "ecommerce-system/internal/services/products"

	"gorm.io/gorm"
)

type CartItemServiceImpl struct {
	*gorm.DB
	cartitemrepositories.CartItemRepositories
	productservices.ProductServices
}

func NewCartItemService(db *gorm.DB, cartrepo cartitemrepositories.CartItemRepositories, productService productservices.ProductServices) CartItemServices {
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
func (cartService *CartItemServiceImpl) handleError(err error) error {
	return exceptions.CheckError(err)
}
func (cartService *CartItemServiceImpl) CreateOrUpdateCartItem(request *request.ReqCreateOrUpdateCartItem, userID int64) (*response.ResCartItem, error) {
	product, err := cartService.ProductServices.GetProductById(request.ProductID)
	if err != nil {
		return nil, err
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

	if errCheck := cartService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	return cartService.loadRes(result), nil
}

func (cartService *CartItemServiceImpl) GetAllUserCartItem(userID int64) ([]*response.ResCartItem, error) {
	results, err := cartService.CartItemRepositories.GetAllUserCartItem(cartService.DB, userID)
	if errCheck := cartService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	items := []*response.ResCartItem{}
	for _, item := range results {
		items = append(items, cartService.loadRes(item))
	}

	return items, nil

}
func (cartService *CartItemServiceImpl) DeleteCartItemById(cartID int64, userID int64) error {

	if err := cartService.CartItemRepositories.DeleteCartItemById(cartService.DB, cartID, userID); err != nil {
		return cartService.handleError(err)
	}

	return nil
}
func (cartService *CartItemServiceImpl) GetCartItemByProductUser(productID int64, userID int64) (*response.ResCartItem, error) {

	result, err := cartService.CartItemRepositories.GetCartItemByProductUser(cartService.DB, productID, userID)
	if errCheck := cartService.handleError(err); errCheck != nil {
		return nil, errCheck
	}

	return cartService.loadRes(result), nil
}
