package checkout

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/domain/model"
	"ecommerce-system/internal/pkg"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

type CheckOutUseCaseImpl struct {
	*gorm.DB
	domain.AddressRepository
	domain.CartRepository
	domain.ProductRepository
	domain.CheckOutRepository
	domain.MidtransGateWay
	domain.PaymentRepository
	domain.OrderRepository
}

func NewCheckOutUseCase(
	coRepo domain.CheckOutRepository,
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	addressRepo domain.AddressRepository,
	paymentRepo domain.PaymentRepository,
	mdGateway domain.MidtransGateWay,
	db *gorm.DB,
) domain.CheckOutUseCase {
	return &CheckOutUseCaseImpl{
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

func (coUC *CheckOutUseCaseImpl) GetLastDraftCheckOut(userID int64) (*response.ResCheckOut, error) {
	result, err := coUC.CheckOutRepository.GetLastDraftCheckOut(coUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

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
		CreatedAt:  result.CreatedAt.Format(pkg.DateTimeLayout),
	}, nil
}

func (coUC *CheckOutUseCaseImpl) CheckOut(req *request.ReqCheckout, userID int64) error {
	var (
		total int64
		items []model.CheckoutItemModel
		err   error
	)
	switch req.Source {
	case "cart":
		total, items, err = coUC.cartItems(req.CartIDs, userID)

	case "direct":
		total, items, err = coUC.directItems(req.Items)

	default:
		return pkg.NewError(pkg.KindBadRequest, "invalid checkout source", nil)
	}

	if err != nil {
		return err
	}

	return coUC.createDraftCheckoutTx(req, items, total, userID)
}

func (coUC *CheckOutUseCaseImpl) CheckOutConfirm(req *request.ReqConfirmCheckout, userID int64) (*response.ResPayment, error) {

	draftCheckout, err := coUC.CheckOutRepository.GetLastDraftCheckOut(coUC.DB, userID)
	if err != nil {
		return nil, pkg.MappingError(err)
	}

	productIDs := make([]int64, 0, len(draftCheckout.CheckoutItem))
	for _, item := range draftCheckout.CheckoutItem {
		productIDs = append(productIDs, item.ProductID)
	}

	sort.Slice(productIDs, func(i, j int) bool {
		return productIDs[i] < productIDs[j]
	})

	orderResult, err := coUC.createOrderTx(req, draftCheckout, productIDs, userID)
	if err != nil {
		return nil, err
	}

	//  MIDTRANS OUTSIDE TRANSACTION
	// KEMUNGKINAN MAKE RABBIT MQ
	snapRes, err := coUC.MidtransGateWay.CreateMidtrans(orderResult)
	if err != nil {
		orderResult.Payment.Status = "failed"
		if errP := coUC.PaymentRepository.SavePayment(coUC.DB, &orderResult.Payment); errP != nil {
			fmt.Println("errUPDATE: ", errP.Error())
		}
		if errSO := coUC.OrderRepository.UpdateStatusOrder(coUC.DB, orderResult.ID, pkg.OderCancel); errSO != nil {
			fmt.Println("errUPDATE: ", errSO.Error())
		}
		if errLC := coUC.CheckOutRepository.UpdateStatusLastCheckOut(coUC.DB, pkg.CheckOutCancel, userID); errLC != nil {
			fmt.Println("errUPDATE: ", errLC.Error())
		}
		return nil, pkg.MappingError(err)
	}

	payment := model.PaymentOrderModel{
		ID:          orderResult.Payment.ID,
		OrderID:     orderResult.ID,
		SnapToken:   snapRes.Token,
		RedirectURL: snapRes.RedirectURL,
		Status:      "pending",
	}

	_ = coUC.PaymentRepository.SavePayment(coUC.DB, &payment)

	return &response.ResPayment{
		ID:          orderResult.Payment.ID,
		OrderID:     orderResult.ID,
		Token:       snapRes.Token,
		RedirectUrl: snapRes.RedirectURL,
	}, nil
}
