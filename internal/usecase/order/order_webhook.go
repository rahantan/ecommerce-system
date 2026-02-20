package order

import (
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"fmt"

	"gorm.io/gorm"
)

func (orderUC *OrderUseCaseImpl) UpdateStatusPayment(transactionOrderID int64, transactionStatus string) error {

	fmt.Println("MASUK USECASE WEBHOOK")
	order, err := orderUC.OrderRepository.GetOrderByID(orderUC.DB, transactionOrderID)
	if err != nil {
		return pkg.MappingError(err)
	}

	if order.StatusID > pkg.OderPending {
		return nil
	}

	switch transactionStatus {
	case pkg.PaymentCancel, pkg.PaymentExpire:
		order.StatusID = pkg.OderCancel

		// update status order jadi cancel
		// update status payment transaksinya sesuai param
		// reverse stock product

		return orderUC.cancelOrderTx(order)

	case pkg.PaymentSettlement:

		// update status order jadi di proses
		// update status payment transaksinya sesuai param
		return orderUC.processOrder(order.ID, pkg.OderProccess, transactionStatus)

	default:
		fmt.Println("STATUS TRANSAKSI WEBHOOK : ", transactionStatus)
		return nil
	}

}

func (orderUC *OrderUseCaseImpl) processOrder(orderID, orderStatusID int64, paymentStatus string) error {
	return orderUC.DB.Transaction(func(tx *gorm.DB) error {
		if err := orderUC.OrderRepository.UpdateStatusOrder(tx, orderID, orderStatusID); err != nil {
			return pkg.MappingError(err)
		}

		if err := orderUC.PaymentRepository.UpdateStatusPayment(tx, orderID, paymentStatus); err != nil {
			return pkg.MappingError(err)
		}

		return nil
	})
}

func (orderUC *OrderUseCaseImpl) cancelOrderTx(order *model.OrderModel) error {
	return orderUC.DB.Transaction(func(tx *gorm.DB) error {

		var (
			productIDs  []int64
			newProducts []model.ProductModel
		)

		for _, item := range order.OrderItem {
			productIDs = append(productIDs, item.ProductID)
		}

		products, err := orderUC.ProductRepository.FindByIDsForUpdate(tx, productIDs)
		if err != nil {
			return pkg.MappingError(err)
		}

		for _, item := range order.OrderItem {
			newProducts = append(newProducts, model.ProductModel{
				ID:    item.ProductID,
				Stock: products[item.ProductID].Stock + item.Qty,
			})
		}

		if err := orderUC.ProductRepository.UpdateProductStock(tx, newProducts); err != nil {
			return pkg.MappingError(err)
		}

		if err := orderUC.OrderRepository.UpdateAll(tx, order); err != nil {
			return pkg.MappingError(err)
		}

		return nil
	})

}
