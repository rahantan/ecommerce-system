package usecase

import (
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

func (orderUC *OrderUseCaseImpl) UpdateStatusPaymentTx(orderID int64, order *model.OrderModel, newPaymentStatus string) error {
	return orderUC.DB.Transaction(func(tx *gorm.DB) error {
		if err := orderUC.PaymentRepository.UpdateStatusPayment(tx, orderID, newPaymentStatus); err != nil {
			return pkg.MappingError(err)
		}

		if order.StatusID == 1 {
			var orderStatusID int64

			switch newPaymentStatus {
			case "settlement":
				orderStatusID = 2
				var (
					productIDs  []int64
					newProducts []model.ProductModel
				)
				for _, item := range order.OrderItem {
					productIDs = append(productIDs, item.ProductID)
				}

				products, err := orderUC.getProductsMap(productIDs)
				if err != nil {
					return err
				}

				for _, item := range order.OrderItem {
					newProducts = append(newProducts, model.ProductModel{
						ID:    item.ProductID,
						Stock: products[item.ProductID].Stock - item.Qty,
					})
				}

				if err = orderUC.ProductRepository.UpdateProductStock(orderUC.DB, newProducts); err != nil {
					return pkg.MappingError(err)
				}

			case "cancel", "expire":
				orderStatusID = 6
			default:
				orderStatusID = 1
			}

			if err := orderUC.OrderRepository.UpdateStatusOrder(tx, orderID, orderStatusID); err != nil {
				return pkg.MappingError(err)
			}

		}

		return nil
	})
}
func (orderUC *OrderUseCaseImpl) UpdateStatusPayment(orderID int64, transactionStatus string) error {

	paymentStatus := map[string]string{
		"settlement": "settlement",
		"cancel":     "cancel",
		"expire":     "expire",
		"refund":     "refund",
		// "pending":        "pending",
		// "capture":        "capture",
		// "deny":           "deny",
		// "partial_refund": "partial_refund",
	}

	payment, err := orderUC.PaymentRepository.GetPaymentByOrderID(orderUC.DB, orderID)
	if err != nil {
		return pkg.MappingError(err)
	}

	if paymentStatus[payment.Status] != "" {
		return nil
	}

	newPaymentStatus := paymentStatus[transactionStatus]

	order, err := orderUC.OrderRepository.GetOrderByID(orderUC.DB, orderID)
	if err != nil {
		return pkg.MappingError(err)
	}

	return orderUC.UpdateStatusPaymentTx(orderID, order, newPaymentStatus)

}
