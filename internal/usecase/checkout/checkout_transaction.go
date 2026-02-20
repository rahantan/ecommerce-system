package checkout

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"

	"gorm.io/gorm"
)

func (coUC *CheckOutUseCaseImpl) createDraftCheckoutTx(req *request.ReqCheckout, items []model.CheckoutItemModel, total int64, userID int64) error {

	return coUC.DB.Transaction(func(tx *gorm.DB) error {

		// CARI TAU ULANG BAGIAN INI UNTUK APA
		if err := coUC.CheckOutRepository.UpdateStatusLastCheckOut(tx, pkg.CheckOutCancel, userID); err != nil {
			return pkg.MappingError(err)
		}

		checkout := model.CheckoutModel{
			Source:       req.Source,
			UserID:       userID,
			TotalPrice:   total,
			Status:       pkg.CheckOutDraft,
			CheckoutItem: items,
		}

		return coUC.CheckOutRepository.CheckOut(tx, &checkout)
	})
}

func (coUC *CheckOutUseCaseImpl) createOrderTx(req *request.ReqConfirmCheckout, draftCheckout *model.CheckoutModel, productIDs []int64, userID int64) (*model.OrderModel, error) {

	var orderResult *model.OrderModel

	err := coUC.DB.Transaction(func(tx *gorm.DB) error {

		products, err := coUC.ProductRepository.FindByIDsForUpdate(tx, productIDs)
		if err != nil {
			return pkg.MappingError(err)
		}

		if err := coUC.validateItems(products, draftCheckout.CheckoutItem); err != nil {
			return err
		}

		var (
			orderItems       []*model.OrderItemModel
			newStockProducts []model.ProductModel
		)

		for _, item := range draftCheckout.CheckoutItem {
			orderItems = append(orderItems, &model.OrderItemModel{
				ProductID: item.ProductID,
				Qty:       item.Qty,
				Price:     item.Price,
				SubTotal:  item.SubTotal,
			})

			newStockProducts = append(newStockProducts, model.ProductModel{
				ID:    item.ProductID,
				Stock: products[item.ProductID].Stock - item.Qty,
			})
		}

		address, err := coUC.resolveAddress(tx, req.AddressID, userID)
		if err != nil {
			return pkg.MappingError(err)
		}

		order := &model.OrderModel{
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
				Status: pkg.PaymentPending,
			},
			UserID: userID,
		}

		orderResult, err = coUC.OrderRepository.CreateOrder(tx, order)
		if err != nil {
			return pkg.MappingError(err)
		}

		if err := coUC.ProductRepository.UpdateProductStock(tx, newStockProducts); err != nil {
			return pkg.MappingError(err)
		}

		if err := coUC.CheckOutRepository.UpdateStatusLastCheckOut(tx, pkg.CheckOutConfirm, userID); err != nil {
			return pkg.MappingError(err)
		}

		if draftCheckout.Source == "cart" {
			if err := coUC.deleteCartItems(tx, draftCheckout.CheckoutItem, userID); err != nil {
				return err
			}
		}

		return nil
	})

	return orderResult, err
}
