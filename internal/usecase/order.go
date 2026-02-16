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

type OrderUseCaseImpl struct {
	*gorm.DB
	domain.CheckOutRepository
	domain.OrderRepository
	domain.CartRepository
	domain.ProductRepository
	domain.AddressRepository
	domain.MidtransGateWay
	domain.PaymentRepository
}

func NewOrderUseCase(
	coRepo domain.CheckOutRepository,
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	addressRepo domain.AddressRepository,
	paymentRepo domain.PaymentRepository,
	mdGateway domain.MidtransGateWay,
	db *gorm.DB,
) domain.OrderUseCase {
	return &OrderUseCaseImpl{
		CheckOutRepository: coRepo,
		OrderRepository:    orderRepo,
		ProductRepository:  productRepo,
		CartRepository:     cartRepo,
		AddressRepository:  addressRepo,
		PaymentRepository:  paymentRepo,
		DB:                 db,
		MidtransGateWay:    mdGateway,
	}
}

func (orderUC *OrderUseCaseImpl) UpdateStatusOrder(orderID, statusOrder int64) error {

	if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, statusOrder); err != nil {
		return pkg.MappingError(err)
	}

	return nil
}

func (orderUC *OrderUseCaseImpl) GetOrderDetails(userID, orderID int64) (*response.ResOrder, error) {
	orderResult, err := orderUC.OrderRepository.GetOrderDetailsByID(orderUC.DB, userID, orderID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	var items []response.ResOrderItem
	for _, item := range orderResult.OrderItem {
		items = append(items, response.ResOrderItem{
			ID:       item.ID,
			Qty:      item.Qty,
			Price:    item.Price,
			SubTotal: item.SubTotal,
			Product: response.ResOrderProduct{
				ID:   item.Product.ID,
				Name: item.Product.Name,
			},
		})
	}
	order := response.ResOrder{
		ID:         orderResult.ID,
		TotalPrice: orderResult.TotalPrice,
		CreatedAt:  orderResult.CreatedAt.Format("2006-01-02 15:04:05"),
		Items:      &items,
		TotalItems: len(orderResult.OrderItem),
		Status: response.ResOrderStatus{
			ID:   orderResult.OrderStatus.ID,
			Name: orderResult.OrderStatus.Name,
		},
	}

	return &order, nil
}
func (orderUC *OrderUseCaseImpl) GetAllOrder(userID int64) ([]response.ResOrder, error) {
	resultOrder, err := orderUC.OrderRepository.GetAllOrder(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	var orders []response.ResOrder

	for _, order := range resultOrder {
		orders = append(orders, response.ResOrder{
			ID: order.ID,
			Status: response.ResOrderStatus{
				ID:   order.OrderStatus.ID,
				Name: order.OrderStatus.Name,
			},
			TotalItems: len(order.OrderItem),
			TotalPrice: order.TotalPrice,
			CreatedAt:  order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return orders, nil
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

func (orderUC *OrderUseCaseImpl) GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error) {
	result, err := orderUC.CheckOutRepository.GetLastDraftCheckOut(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}
	fmt.Println("ITEMS: ", result.CheckoutItem)

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

////////////////////////////////////////////////////////////
//////////////////// CHECKOUT FLOW /////////////////////////
////////////////////////////////////////////////////////////

func (orderUC *OrderUseCaseImpl) CheckOut(req *request.ReqCheckout, userID int64) error {
	var (
		total int64
		items []model.CheckoutItemModel
		err   error
	)
	switch req.Source {
	case "cart":
		total, items, err = orderUC.cartItems(req.CartIDs, userID)

	case "direct":
		total, items, err = orderUC.directItems(req.Items)

	default:
		return pkg.NewError(pkg.KindBadRequest, "invalid checkout source", nil)
	}

	if err != nil {
		return err
	}

	return orderUC.createDraftCheckoutTx(req, items, total, userID)
}

func (orderUC *OrderUseCaseImpl) CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error) {

	draftCheckout, err := orderUC.CheckOutRepository.GetLastDraftCheckOut(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	productIDs := make([]int64, 0, len(draftCheckout.CheckoutItem))
	for _, item := range draftCheckout.CheckoutItem {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := orderUC.getProductsMap(productIDs)
	if err != nil {
		return nil, err
	}

	if err := orderUC.validateItems(products, draftCheckout.CheckoutItem); err != nil {
		return nil, err
	}

	orderResult, err := orderUC.createOrderTx(req, draftCheckout, userID)
	if err != nil {
		return nil, err
	}

	//  MIDTRANS OUTSIDE TRANSACTION
	snapRes, err := orderUC.MidtransGateWay.CreateMidtrans(orderResult)
	if err != nil {
		fmt.Println("payment: ", orderResult.Payment)
		orderResult.Payment.Status = "failed"
		fmt.Println("payment: ", orderResult.Payment)
		if errP := orderUC.PaymentRepository.SavePayment(orderUC.DB, &orderResult.Payment); errP != nil {
			fmt.Println("errUPDATE: ", errP.Error())
		}
		if errSO := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderResult.ID, 6); errSO != nil {
			fmt.Println("errUPDATE: ", errSO.Error())
		}
		if errLC := orderUC.CheckOutRepository.UpdateStatusLastCheckOut(orderUC.DB, "cancel", userID); errLC != nil {
			fmt.Println("errUPDATE: ", errLC.Error())
		}
		return nil, pkg.MappingError(err)
	}

	payment := model.PaymentOrderModel{
		OrderID:     orderResult.ID,
		SnapToken:   snapRes.Token,
		RedirectURL: snapRes.RedirectURL,
		Status:      "pending",
	}

	fmt.Println(" PAYMENT: ", payment)
	_ = orderUC.PaymentRepository.SavePayment(orderUC.DB, &payment)

	return &response.ResPayment{
		OrderID:     orderResult.ID,
		Token:       snapRes.Token,
		RedirectUrl: snapRes.RedirectURL,
	}, nil
}

////////////////////////////////////////////////////////////
//////////////////// TRANSACTIONS //////////////////////////
////////////////////////////////////////////////////////////

func (orderUC *OrderUseCaseImpl) createDraftCheckoutTx(req *request.ReqCheckout, items []model.CheckoutItemModel, total int64, userID int64) error {

	return orderUC.DB.Transaction(func(tx *gorm.DB) error {

		if err := orderUC.CheckOutRepository.
			UpdateStatusLastCheckOut(tx, "cancel", userID); err != nil {
			return pkg.MappingError(err)
		}

		checkout := model.CheckoutModel{
			Source:       req.Source,
			UserID:       userID,
			TotalPrice:   total,
			Status:       "draft",
			CheckoutItem: items,
		}

		return orderUC.CheckOutRepository.CheckOut(tx, &checkout)
	})
}

func (orderUC *OrderUseCaseImpl) createOrderTx(req *request.ReqConfirmCheckout, draftCheckout *model.CheckoutModel, userID int64) (*model.OrderModel, error) {

	var orderResult *model.OrderModel

	err := orderUC.DB.Transaction(func(tx *gorm.DB) error {

		var orderItems []*model.OrderItemModel

		for _, item := range draftCheckout.CheckoutItem {

			orderItems = append(orderItems, &model.OrderItemModel{
				ProductID: item.ProductID,
				Qty:       item.Qty,
				Price:     item.Price,
				SubTotal:  item.SubTotal,
			})
		}

		address, err := orderUC.resolveAddress(tx, req.AddressID, userID)
		if err != nil {
			return pkg.MappingError(err)
		}

		order := model.OrderModel{
			TotalPrice:    draftCheckout.TotalPrice,
			PaymentMethod: req.PaymentMethod,
			Noted:         req.Note,
			StatusID:      1,
			OrderItem:     orderItems,
			AddressOrder: model.AddressOrderModel{
				City:    address.City,
				Address: address.Address,
			},
			Payment: model.PaymentOrderModel{
				Status: "pending",
			},
			UserID: userID,
		}

		result, err := orderUC.OrderRepository.CreateOrder(tx, &order)
		if err != nil {
			return pkg.MappingError(err)
		}

		if err := orderUC.CheckOutRepository.UpdateStatusLastCheckOut(tx, "confirm", userID); err != nil {
			return pkg.MappingError(err)
		}

		if draftCheckout.Source == "cart" {
			if err := orderUC.deleteCartItems(tx, draftCheckout.CheckoutItem, userID); err != nil {
				return err
			}
		}
		orderResult = result

		return nil
	})

	return orderResult, err
}

// HELPERS

func (orderUC *OrderUseCaseImpl) resolveAddress(tx *gorm.DB, addressID int64, userID int64) (*model.AddressModel, error) {

	if addressID == 0 {
		return orderUC.AddressRepository.GetUserAddressActive(tx, userID)
	}

	return orderUC.AddressRepository.GetAddressById(tx, addressID)
}

func (orderUC *OrderUseCaseImpl) deleteCartItems(tx *gorm.DB, items []model.CheckoutItemModel, userID int64) error {

	var cartIDs []int64
	for _, item := range items {
		if item.CartID != nil {
			cartIDs = append(cartIDs, *item.CartID)
		}
	}

	err := orderUC.CartRepository.DeleteCartItemsByIDs(tx, cartIDs, userID)
	if err != nil {
		return pkg.MappingError(err)
	}

	return nil

}

func (orderUC *OrderUseCaseImpl) getProductsMap(productIDs []int64) (map[int64]*model.ProductModel, error) {

	products, err := orderUC.ProductRepository.GetAllProductByIDs(orderUC.DB, productIDs)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	result := make(map[int64]*model.ProductModel)
	for _, p := range products {
		result[p.ID] = p
	}

	return result, nil
}

func (orderUC *OrderUseCaseImpl) validateItems(products map[int64]*model.ProductModel, items []model.CheckoutItemModel) error {

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

func (orderUC *OrderUseCaseImpl) cartItems(cartIDs []int64, userID int64) (int64, []model.CheckoutItemModel, error) {

	items, err := orderUC.CartRepository.GetAllCartItemByIDs(orderUC.DB, cartIDs, userID)
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

		if item.Qty > item.Product.Stock {
			return 0, nil, pkg.NewError(pkg.KindConflict, "qty cannot exceed available stock", nil)
		}

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
func (orderUC *OrderUseCaseImpl) directItems(reqItems []request.ReqItem) (int64, []model.CheckoutItemModel, error) {

	var (
		checkOutItems []model.CheckoutItemModel
		productIDs    []int64
		totalPrice    int64
	)

	for _, item := range reqItems {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := orderUC.getProductsMap(productIDs)
	if err != nil {
		return 0, nil, err
	}

	for _, item := range reqItems {
		if item.Qty > products[item.ProductID].Stock {
			return 0, nil, pkg.NewError(pkg.KindConflict, "qty cannot exceed available stock", nil)
		}

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
