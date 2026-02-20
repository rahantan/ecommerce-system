package order

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
)

func (orderUC *OrderUseCaseImpl) resOrder(order *model.OrderModel, items *[]response.ResOrderItem) *response.ResOrder {
	return &response.ResOrder{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format(pkg.DateTimeLayout),
		Items:      items,
		TotalItems: len(order.OrderItem),
		Status: response.ResOrderStatus{
			ID:   order.OrderStatus.ID,
			Name: order.OrderStatus.Name,
		},
	}
}

func (orderUC *OrderUseCaseImpl) resOrderItem(item *model.OrderItemModel) response.ResOrderItem {
	return response.ResOrderItem{
		ID:       item.ID,
		Qty:      item.Qty,
		Price:    item.Price,
		SubTotal: item.SubTotal,
		Product: response.ResOrderProduct{
			ID:   item.Product.ID,
			Name: item.Product.Name,
		},
	}
}
