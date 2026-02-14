package midtrans

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"errors"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransGateWayImpl struct {
	snap.Client
}

func NewMidtransGateWay(mdClient snap.Client) domain.MidtransGateWay {
	return &MidtransGateWayImpl{
		Client: mdClient,
	}
}
func (md *MidtransGateWayImpl) CreateMidtrans(order *model.OrderModel) (*snap.Response, error) {

	if len(order.OrderItem) < 1 {
		return nil, pkg.NewError(pkg.KindCancelled, "no items to create midtrans", nil)
	}

	var items []midtrans.ItemDetails
	for _, item := range order.OrderItem {

		fmt.Println("product id", item.ProductID)
		fmt.Println("product name", item.Product.Name)
		fmt.Println("product price", item.Price)
		fmt.Println("qty", item.Qty)
		items = append(items, midtrans.ItemDetails{
			ID:    strconv.Itoa(int(item.ProductID)),
			Name:  item.Product.Name,
			Price: int64(item.Price),
			Qty:   int32(item.Qty),
		})
	}

	orderID := order.ID
	fmt.Println("saat confirm checkout", orderID)
	snapRes, err := md.Client.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(order.ID)),
			GrossAmt: int64(order.TotalPrice),
		},
		Items: &items,
	})
	if err != nil {
		return nil, errors.New(err.Message)
	}
	return snapRes, nil
}
