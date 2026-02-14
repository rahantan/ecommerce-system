package usecase

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"fmt"
	"net/http"

	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type OrderUseCaseImpl struct {
	*gorm.DB
	domain.CheckOutRepository
	domain.OrderRepository
	domain.CartRepository
	domain.ProductRepository
	domain.AddressRepository
	domain.MidtransGateWay
}

func NewOrderUseCase(
	coRepo domain.CheckOutRepository,
	orderRepo domain.OrderRepository,
	cartrepo domain.CartRepository,
	productRepo domain.ProductRepository,
	addressRepo domain.AddressRepository,
	mdGateWay domain.MidtransGateWay,
	db *gorm.DB,
) domain.OrderUseCase {
	return &OrderUseCaseImpl{
		CheckOutRepository: coRepo,
		OrderRepository:    orderRepo,
		ProductRepository:  productRepo,
		CartRepository:     cartrepo,
		AddressRepository:  addressRepo,
		DB:                 db,
		MidtransGateWay:    mdGateWay,
	}
}
func (orderUC *OrderUseCaseImpl) validationProduct(products map[int64]*model.ProductModel, checkOutItems []model.CheckoutItemModel) (float64, error) {

	var totalPrice float64 = 0
	for _, item := range checkOutItems {
		//product not exist
		if products[item.ProductID].ID != item.ProductID {
			return 0, pkg.NewError(pkg.KindNotFound, fmt.Sprintf("product id %d not found", item.ProductID), http.StatusNotFound)
		}

		//stock not exist
		if item.Qty > products[item.ProductID].Stock {
			return 0, pkg.NewError(pkg.KindCancelled, "qty cannot exceed available stock", http.StatusBadRequest)
		}
		totalPrice += item.SubTotal
	}
	return totalPrice, nil
}

func (orderUC *OrderUseCaseImpl) getItemByCarts(products map[int64]*model.ProductModel, cartIDs []int64, userID int64) ([]model.CheckoutItemModel, error) {
	checkOutItems := []model.CheckoutItemModel{}

	items, err := orderUC.CartRepository.GetAllCartItemByIDs(orderUC.DB, cartIDs, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	if len(cartIDs) != len(items) {
		return nil, pkg.NewError(pkg.KindNotFound, "cart item id not exist, please check your cart id", nil)
	}

	for _, item := range items {
		checkOutItems = append(checkOutItems, model.CheckoutItemModel{
			CartID:    &item.ID,
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     products[item.ProductID].Price,
			SubTotal:  products[item.ProductID].Price * float64(item.Qty),
		})
	}

	return checkOutItems, nil
}
func (orderUC *OrderUseCaseImpl) mappingItems(products map[int64]*model.ProductModel, items []request.ReqItem) ([]model.CheckoutItemModel, error) {

	checkOutItems := []model.CheckoutItemModel{}
	for _, item := range items {
		checkOutItems = append(checkOutItems, model.CheckoutItemModel{
			ProductID: item.ProductID,
			Qty:       item.Qty,
			Price:     products[item.ProductID].Price,
			SubTotal:  products[item.ProductID].Price * float64(item.Qty),
		})
	}

	return checkOutItems, nil
}

func (orderUC *OrderUseCaseImpl) getProductsWithMapping() (map[int64]*model.ProductModel, error) {

	getProducts, err := orderUC.ProductRepository.GetAllProduct(orderUC.DB)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	products := make(map[int64]*model.ProductModel)
	for _, product := range getProducts {
		products[product.ID] = product
	}

	return products, nil
}
func (orderUC *OrderUseCaseImpl) checkOutTx(req *request.ReqCheckout, checkOutItems []model.CheckoutItemModel, totalPrice float64, userID int64) error {
	return orderUC.DB.Transaction(func(tx *gorm.DB) error {

		if err := orderUC.CheckOutRepository.UpdateStatusLastCheckOut(tx, "cancel", userID); err != nil {
			return pkg.MappingError(err)
		}

		checkOut := model.CheckoutModel{
			Source:       req.Source,
			UserID:       userID,
			TotalPrice:   totalPrice,
			Status:       "draft",
			CheckoutItem: checkOutItems,
		}

		if err := orderUC.CheckOutRepository.CheckOut(tx, &checkOut); err != nil {
			fmt.Println("error update cekout")
			return pkg.MappingError(err)
		}

		return nil
	})
}

func (orderUC *OrderUseCaseImpl) GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error) {
	result, err := orderUC.CheckOutRepository.GetLastDraftCheckOut(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
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

func (orderUC *OrderUseCaseImpl) CheckOut(req *request.ReqCheckout, userID int64) error {

	products, err := orderUC.getProductsWithMapping()
	if err != nil {
		return pkg.MappingError(err)
	}

	var checkOutItems []model.CheckoutItemModel

	switch req.Source {
	case "cart":
		items, err := orderUC.getItemByCarts(products, req.CartIDs, userID)
		if err != nil {
			return err
		}

		checkOutItems = items

	case "direct":
		items, err := orderUC.mappingItems(products, req.Items)
		if err != nil {
			return err
		}
		checkOutItems = items
	}

	totalPrice, err := orderUC.validationProduct(products, checkOutItems)
	if err != nil {
		return err
	}

	if err := orderUC.checkOutTx(req, checkOutItems, totalPrice, userID); err != nil {
		return err
	}

	return nil
}
func (orderUC *OrderUseCaseImpl) CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error) {

	products, err := orderUC.getProductsWithMapping()
	if err != nil {
		return nil, err
	}

	draftCheckOut, err := orderUC.CheckOutRepository.GetLastDraftCheckOut(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	_, err = orderUC.validationProduct(products, draftCheckOut.CheckoutItem)
	if err != nil {
		return nil, err
	}

	snapRes, err := orderUC.checkOutConfirmTx(req, draftCheckOut, products, userID)
	if err != nil {
		return nil, err
	}
	return &response.ResPayment{
		Token:       snapRes.Token,
		RedirectUrl: snapRes.RedirectURL,
	}, nil
}

func (orderUC *OrderUseCaseImpl) checkOutConfirmTx(req *request.ReqConfirmCheckout, draftCheckOut *model.CheckoutModel, products map[int64]*model.ProductModel, userID int64) (*snap.Response, error) {
	var snapRes *snap.Response
	err := orderUC.DB.Transaction(func(tx *gorm.DB) error {

		productStockUpdate := []*model.ProductModel{}

		orderItems := []*model.OrderItemModel{}
		for _, item := range draftCheckOut.CheckoutItem {
			orderItems = append(orderItems, &model.OrderItemModel{
				ProductID: item.ProductID,
				Qty:       item.Qty,
				Price:     item.Price,
				SubTotal:  item.SubTotal,
			})

			productStockUpdate = append(productStockUpdate, &model.ProductModel{
				ID:    item.ProductID,
				Stock: products[item.ProductID].Stock - item.Qty,
			})
		}

		var address *model.AddressModel
		if req.AddressID == 0 {
			result, err := orderUC.AddressRepository.GetUserAddressActive(tx, userID)
			if err != nil {
				return pkg.MappingError(err)
			}
			address = result
		} else {
			result, err := orderUC.AddressRepository.GetAddressById(tx, req.AddressID)
			if err != nil {
				return pkg.MappingError(err)
			}
			address = result
		}

		orderCO := model.OrderModel{
			TotalPrice:    draftCheckOut.TotalPrice,
			PaymentMethod: req.PaymentMethod,
			Noted:         req.Note,
			StatusID:      1,
			OrderItem:     orderItems,
			AddressOrder: model.AddressOrderModel{
				City:    address.City,
				Address: address.Address,
			},
			UserID: userID,
		}

		orderResult, err := orderUC.OrderRepository.CreateOrder(tx, &orderCO)
		if err != nil {
			return pkg.MappingError(err)
		}

		if err := orderUC.CheckOutRepository.UpdateStatusLastCheckOut(tx, "confirm", userID); err != nil {
			return pkg.MappingError(err)
		}

		//delete items cart
		if draftCheckOut.Source == "cart" {
			cartIDs := []int64{}
			for _, item := range draftCheckOut.CheckoutItem {
				cartIDs = append(cartIDs, *item.CartID)
			}

			if err := orderUC.CartRepository.DeleteCartItemsByIDs(tx, cartIDs, userID); err != nil {
				return pkg.MappingError(err)
			}
		}
		if err := orderUC.ProductRepository.UpdateProductStockByID(tx, productStockUpdate); err != nil {
			return pkg.MappingError(err)
		}

		snap, err := orderUC.MidtransGateWay.CreateMidtrans(orderResult)
		if err != nil {
			return pkg.NewError(pkg.KindInternal, err.Error(), nil)
		}
		snapRes = snap
		return nil
	})
	if err != nil {
		return nil, err
	}
	return snapRes, nil
}
func (orderUC *OrderUseCaseImpl) UpdateStatusOrder(orderID, statusOrder int64) error {

	order, err := orderUC.OrderRepository.GetOrderByID(orderUC.DB, orderID)
	if err != nil {
		return pkg.MappingError(err)
	}
	if order.StatusID > 1 {
		return nil
	}
	if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, statusOrder); err != nil {
		return pkg.MappingError(err)
	}
	return nil

}
