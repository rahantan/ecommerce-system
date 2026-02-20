package order

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"

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

func (orderUC *OrderUseCaseImpl) ShipOrder(orderID int64) error {
	order, err := orderUC.OrderRepository.GetOrderByID(orderUC.DB, orderID)
	if err != nil {
		return pkg.MappingError(err)
	}
	switch order.StatusID {

	case pkg.OderPending:
		return pkg.NewError(pkg.KindCancelled, "Order has not been paid yet; cannot update", nil)

	case pkg.OderShip:
		return pkg.NewError(pkg.KindInfo, "Order already shipped; no update performed", nil)

	case pkg.OderReceive:
		return pkg.NewError(pkg.KindInfo, "Order already received; no update performed", nil)

	case pkg.OderCancel:
		return pkg.NewError(pkg.KindCancelled, "Order has been cancelled; cannot update", nil)

	default:
		if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, pkg.OderShip); err != nil {
			return pkg.MappingError(err)
		}
	}

	return nil

}

func (orderUC *OrderUseCaseImpl) ReceiveOrder(orderID, userID int64) error {
	order, err := orderUC.OrderRepository.GetOrderByID(orderUC.DB, orderID)
	if err != nil {
		return pkg.MappingError(err)
	}
	switch order.StatusID {

	case pkg.OderPending:
		return pkg.NewError(pkg.KindCancelled, "Order has not been paid yet; cannot update", nil)

	case pkg.OderProccess:
		return pkg.NewError(pkg.KindCancelled, "Order is still being processed; cannot update", nil)

	case pkg.OderReceive:
		return pkg.NewError(pkg.KindInfo, "Order already received; no update performed", nil)

	case pkg.OderCancel:
		return pkg.NewError(pkg.KindCancelled, "Order has been cancelled; cannot update", nil)

	default:
		if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, pkg.OderReceive); err != nil {
			return pkg.MappingError(err)
		}
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
		items = append(items, orderUC.resOrderItem(item))
	}

	return orderUC.resOrder(orderResult, &items), nil
}

func (orderUC *OrderUseCaseImpl) GetAllOrder() ([]*response.ResOrder, error) {
	resultOrder, err := orderUC.OrderRepository.GetAllOrder(orderUC.DB)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	var orders []*response.ResOrder

	for _, order := range resultOrder {
		orders = append(orders, orderUC.resOrder(order, nil))
	}

	return orders, nil
}

func (orderUC *OrderUseCaseImpl) GetAllOrderByUserID(userID int64) ([]*response.ResOrder, error) {
	resultOrder, err := orderUC.OrderRepository.GetAllOrderByUserID(orderUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	var orders []*response.ResOrder

	for _, order := range resultOrder {
		orders = append(orders, orderUC.resOrder(order, nil))
	}

	return orders, nil
}
