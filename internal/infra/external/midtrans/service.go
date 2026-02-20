package midtrans

import (
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"errors"
	"strconv"
	"time"

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

		// item.Product.Name = ""
		if item.Product.Name == "" {
			return nil, pkg.NewError(pkg.KindNotFound, "product name not found", nil)
		}

		items = append(items, midtrans.ItemDetails{
			ID:    strconv.Itoa(int(item.ProductID)),
			Name:  item.Product.Name,
			Price: int64(item.Price),
			Qty:   int32(item.Qty),
		})
	}

	snapRes, err := md.Client.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(order.ID)),
			GrossAmt: int64(order.TotalPrice),
		},
		Items: &items,
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format(pkg.DateTimeLayout + " -0700"),
			Unit:      "minute",
			Duration:  10,
		},
	})
	if err != nil {
		return nil, errors.New(err.Message)
	}
	return snapRes, nil
}
