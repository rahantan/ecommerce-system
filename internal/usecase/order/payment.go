package order

import (
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/pkg"
)

func (orderUC *OrderUseCaseImpl) GetUserPaymentByOrderID(orderID, userID int64) (*response.ResPayment, error) {
	payment, err := orderUC.PaymentRepository.GetPaymentByOrderID(orderUC.DB, orderID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}
	disallowedPaymentStatuses := map[string]string{
		pkg.PaymentSettlement: "The payment has been successfully completed.",
		pkg.PaymentCancel:     "The payment has been canceled.",
		pkg.PaymentExpire:     "The payment has expired.",
		pkg.PaymentRefund:     "The payment has been refunded.",
	}

	if disallowedPaymentStatuses[payment.Status] != "" {
		return nil, pkg.NewError(pkg.KindInfo, disallowedPaymentStatuses[payment.Status], nil)
	}

	return &response.ResPayment{
		ID:          payment.ID,
		Token:       payment.SnapToken,
		OrderID:     payment.OrderID,
		RedirectUrl: payment.RedirectURL,
	}, nil

}
