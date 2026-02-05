package orderservices

import (
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	addressrepositories "ecommerce-system/internal/repositories/addresses"
	cartitemrepositories "ecommerce-system/internal/repositories/carts"
	orderrepositories "ecommerce-system/internal/repositories/orders"
	productrepositories "ecommerce-system/internal/repositories/products"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type OrderServiceImpl struct {
	*gorm.DB
	orderrepositories.CheckOutRepositories
	orderrepositories.OrderRepositories
	cartitemrepositories.CartItemRepositories
	productrepositories.ProductRepositories
	addressrepositories.AddressRepositories
}

func NewOrderService(
	coRepo orderrepositories.CheckOutRepositories,
	orderRepo orderrepositories.OrderRepositories,
	cartrepo cartitemrepositories.CartItemRepositories,
	productRepo productrepositories.ProductRepositories,
	addressRepo addressrepositories.AddressRepositories,
	db *gorm.DB,
) OrderServices {
	return &OrderServiceImpl{
		CheckOutRepositories: coRepo,
		OrderRepositories:    orderRepo,
		ProductRepositories:  productRepo,
		CartItemRepositories: cartrepo,
		AddressRepositories:  addressRepo,
		DB:                   db,
	}
}
func (coService *OrderServiceImpl) handleError(err error) error {
	newErr := exceptions.CheckError(err)
	return newErr
}

func (coService *OrderServiceImpl) getItemByCarts(products map[int64]*models.ProductModel, cartIDs []int64, userID int64) ([]models.CheckoutItemModel, error) {
	checkOutItems := []models.CheckoutItemModel{}

	items, err := coService.CartItemRepositories.GetAllCartItemByIDs(coService.DB, cartIDs, userID)
	if err != nil {
		return nil, err
	}
	if len(cartIDs) != len(items) {
		return nil, exceptions.NewError("", "cart item id not exist, please check your cart id", http.StatusBadRequest)
	}

	for _, item := range items {
		checkOutItems = append(checkOutItems, models.CheckoutItemModel{
			CartID:    item.ID,
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     products[item.ProductID].Price,
			SubTotal:  products[item.ProductID].Price * float64(item.Qty),
		})
	}

	return checkOutItems, nil
}
func (coService *OrderServiceImpl) mappingItems(products map[int64]*models.ProductModel, items []request.ReqItem) ([]models.CheckoutItemModel, error) {

	checkOutItems := []models.CheckoutItemModel{}
	for _, item := range items {
		checkOutItems = append(checkOutItems, models.CheckoutItemModel{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     products[item.ProductID].Price,
			SubTotal:  products[item.ProductID].Price * float64(item.Qty),
		})
	}

	return checkOutItems, nil
}
func (coService *OrderServiceImpl) validationProduct(products map[int64]*models.ProductModel, checkOutItems []models.CheckoutItemModel) (float64, error) {
	//validate product
	var totalPrice float64 = 0
	for _, item := range checkOutItems {
		//product exist
		if products[item.ProductID].ID != item.ProductID {
			return 0, coService.handleError(exceptions.NewError("", fmt.Sprintf("product id %d not found", item.ProductID), http.StatusNotFound))
		}
		//stock exist
		if item.Qty > products[item.ProductID].Stock {
			return 0, coService.handleError(exceptions.NewError("", "qty cannot exceed available stock", http.StatusBadRequest))
		}
		totalPrice += item.SubTotal
	}
	return totalPrice, nil
}
func (coService *OrderServiceImpl) getProductsWithMapping() (map[int64]*models.ProductModel, error) {

	getProducts, err := coService.ProductRepositories.GetAllProduct(coService.DB)
	if err != nil {
		return nil, err
	}

	products := make(map[int64]*models.ProductModel)
	for _, product := range getProducts {
		products[product.ID] = product
	}

	return products, nil
}
func (coService *OrderServiceImpl) checkOutTx(req *request.ReqCheckout, checkOutItems []models.CheckoutItemModel, totalPrice float64, userID int64) error {
	return coService.DB.Transaction(func(tx *gorm.DB) error {
		fmt.Println(0)
		if err := coService.CheckOutRepositories.UpdateStatusLastCheckOut(tx, "cancel", userID); err != nil {
			return err
		}
		fmt.Println(1)

		checkOut := models.CheckoutModel{
			Source:       req.Source,
			UserID:       userID,
			TotalPrice:   totalPrice,
			Status:       "draft",
			CheckoutItem: checkOutItems,
		}

		if err := coService.CheckOutRepositories.CheckOut(tx, &checkOut); err != nil {

			return err
		}
		return nil
	})
}
func (coService *OrderServiceImpl) checkOutConfirmTx(products map[int64]*models.ProductModel, req *request.ReqConfirmCheckout, draftCheckOut *models.CheckoutModel, userID int64) error {
	return coService.DB.Transaction(func(tx *gorm.DB) error {

		productStockUpdate := []*models.ProductModel{}

		orderItems := []*models.OrderItemModel{}
		for _, item := range draftCheckOut.CheckoutItem {
			orderItems = append(orderItems, &models.OrderItemModel{
				ProductID: item.ProductID,
				Qty:       item.Qty,
				Price:     item.Price,
				SubTotal:  item.SubTotal,
			})

			productStockUpdate = append(productStockUpdate, &models.ProductModel{
				ID:    item.ProductID,
				Stock: products[item.ProductID].Stock - item.Qty,
			})
		}

		var address *models.AddressModel
		if req.AddressID == 0 {
			result, err := coService.AddressRepositories.GetUserAddressActive(tx, userID)
			if err != nil {
				return err
			}
			address = result
		} else {
			result, err := coService.AddressRepositories.GetAddressById(tx, req.AddressID)
			if err != nil {
				return err
			}
			address = result
		}

		addressOrder := models.AddressOrderModel{
			City:    address.City,
			Address: address.Address,
		}

		order := models.OrderModel{
			TotalPrice:    draftCheckOut.TotalPrice,
			PaymentMethod: req.PaymentMethod,
			Noted:         req.Note,
			StatusID:      1,
			OrderItem:     orderItems,
			AddressOrder:  addressOrder,
			UserID:        userID,
		}

		if err := coService.OrderRepositories.CreateOrder(tx, &order); err != nil {
			return err
		}

		if err := coService.CheckOutRepositories.UpdateStatusLastCheckOut(tx, "confirm", userID); err != nil {
			return err
		}

		if draftCheckOut.Source == "cart" {
			cartIDs := []int64{}
			for _, item := range draftCheckOut.CheckoutItem {
				cartIDs = append(cartIDs, item.CartID)
			}

			if err := coService.CartItemRepositories.DeleteCartItemsByIDs(tx, cartIDs, userID); err != nil {
				return err
			}
		}
		if err := coService.ProductRepositories.UpdateProductStockByID(tx, productStockUpdate); err != nil {
			return err
		}

		return nil
	})
}
func (coService *OrderServiceImpl) GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error) {
	result, err := coService.CheckOutRepositories.GetLastDraftCheckOut(coService.DB, userID)
	if err != nil {
		return nil, coService.handleError(err)
	}

	var items []response.ResItem
	for _, item := range result.CheckoutItem {
		items = append(items, response.ResItem{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     item.Price,
			SubTotal:  item.SubTotal,
		})
	}

	return &response.ResCheckOut{
		ID:         result.ID,
		Status:     result.Status,
		Items:      items,
		TotalPrice: result.TotalPrice,
		CreatedAt:  result.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (coService *OrderServiceImpl) CheckOut(req *request.ReqCheckout, userID int64) error {

	products, err := coService.getProductsWithMapping()
	if err != nil {
		return coService.handleError(err)
	}

	var checkOutItems []models.CheckoutItemModel

	switch req.Source {
	case "cart":

		items, err := coService.getItemByCarts(products, req.CartIDs, userID)
		if err != nil {
			return coService.handleError(err)
		}

		checkOutItems = items

	case "direct":
		items, err := coService.mappingItems(products, req.Items)
		if err != nil {
			return coService.handleError(err)
		}
		checkOutItems = items
	}

	totalPrice, err := coService.validationProduct(products, checkOutItems)
	if err != nil {
		return coService.handleError(err)
	}

	if err := coService.checkOutTx(req, checkOutItems, totalPrice, userID); err != nil {
		return coService.handleError(err)
	}

	return nil
}
func (coService *OrderServiceImpl) CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) error {

	products, err := coService.getProductsWithMapping()
	if err != nil {
		return coService.handleError(err)
	}

	draftCheckOut, err := coService.CheckOutRepositories.GetLastDraftCheckOut(coService.DB, userID)
	if err != nil {
		return coService.handleError(err)
	}

	_, err = coService.validationProduct(products, draftCheckOut.CheckoutItem)
	if err != nil {
		return coService.handleError(err)
	}

	if err := coService.checkOutConfirmTx(products, req, draftCheckOut, userID); err != nil {
		return coService.handleError(err)
	}

	return nil
}
