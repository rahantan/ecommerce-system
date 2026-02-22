package order

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"
	"math"

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

	case pkg.OrderPending:
		return pkg.NewError(pkg.KindCancelled, "Order has not been paid yet; cannot update", nil)

	case pkg.OrderShip:
		return pkg.NewError(pkg.KindInfo, "Order already shipped; no update performed", nil)

	case pkg.OrderReceive:
		return pkg.NewError(pkg.KindInfo, "Order already received; no update performed", nil)

	case pkg.OrderCancel:
		return pkg.NewError(pkg.KindCancelled, "Order has been cancelled; cannot update", nil)

	default:
		if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, pkg.OrderShip); err != nil {
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

	case pkg.OrderPending:
		return pkg.NewError(pkg.KindCancelled, "Order has not been paid yet; cannot update", nil)

	case pkg.OrderProccess:
		return pkg.NewError(pkg.KindCancelled, "Order is still being processed; cannot update", nil)

	case pkg.OrderReceive:
		return pkg.NewError(pkg.KindInfo, "Order already received; no update performed", nil)

	case pkg.OrderCancel:
		return pkg.NewError(pkg.KindCancelled, "Order has been cancelled; cannot update", nil)

	default:
		if err := orderUC.OrderRepository.UpdateStatusOrder(orderUC.DB, orderID, pkg.OrderReceive); err != nil {
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

func (orderUC *OrderUseCaseImpl) GetAllOrderByUserID(page, limit int, userID int64) ([]*response.ResOrder, *response.ResPaginateStandard, error) {
	resultOrder, count, err := orderUC.OrderRepository.GetAllOrderByUserID(orderUC.DB, page, limit, userID)
	if err != nil {
		return nil, nil, pkg.MappingError(err)
	}

	var orders []*response.ResOrder

	for _, order := range resultOrder {
		orders = append(orders, orderUC.resOrder(order, nil))
	}
	totalPage := math.Ceil(float64(count) / float64(limit))

	paginate := &response.ResPaginateStandard{
		Page:      page,
		Limit:     limit,
		TotalData: int(count),
		TotalPage: int(totalPage),
	}

	return orders, paginate, nil
}
