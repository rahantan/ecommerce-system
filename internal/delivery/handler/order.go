package handler

import (
	"crypto/sha512"
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandlerImpl struct {
	domain.OrderUseCase
	config.Midtrans
	*validator.Validate
}

func NewOrderHandler(order domain.OrderUseCase, md config.Midtrans, v *validator.Validate) domain.OrderHandler {
	return &OrderHandlerImpl{
		OrderUseCase: order,
		Midtrans:     md,
		Validate:     v,
	}
}
func (orderHandler *OrderHandlerImpl) midTransVerify(orderId, statusCode, grossMount, signature string) bool {
	payLoad := orderId + statusCode + grossMount + orderHandler.Midtrans.ServerKey
	hash := sha512.Sum512([]byte(payLoad))
	return hex.EncodeToString(hash[:]) == signature
}
func (orderHandler *OrderHandlerImpl) WebHookMidtransNotif(ctx *fiber.Ctx) error {
	var body request.MidtransNotification
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if !orderHandler.midTransVerify(body.OrderID, body.StatusCode, body.GrossAmount, body.SignatureKey) {
		return ctx.SendStatus(pkg.ErrCustomForbidden.GetStatusCode())
	}

	if body.FraudStatus == "accept" {
		orderID, err := strconv.Atoi(body.OrderID)
		if err != nil {
			return pkg.ErrCustomInvalidOrderId
		}
		fmt.Println("order id: ", orderID)

		if err := orderHandler.UpdateStatusPayment(int64(orderID), body.TransactionStatus); err != nil {
			return err
		}

	} else {

		fmt.Println("ente mau nipu yah?")
		return ctx.SendStatus(pkg.ErrCustomForbidden.GetStatusCode())
	}
	return ctx.SendStatus(200)
}
func (orderHandler *OrderHandlerImpl) getUserData(ctx *fiber.Ctx) (*response.ResUser, error) {
	user, ok := ctx.Locals("user").(response.ResUser)
	if !ok {
		return nil, pkg.ErrCustomUnauthorized
	}

	return &user, nil
}

func (orderHandler *OrderHandlerImpl) GetUserPaymentByOrderID(ctx *fiber.Ctx) error {
	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	result, err := orderHandler.OrderUseCase.GetUserPaymentByOrderID(int64(orderID), user.ID)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success confirm checkout",
		Data: map[string]any{
			"payment": result,
		},
	})
}
func (orderHandler *OrderHandlerImpl) ReceiveOrder(ctx *fiber.Ctx) error {
	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	if err := orderHandler.OrderUseCase.ReceiveOrder(int64(orderID), user.ID); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "Order successfully marked as received.",
	})
}
func (orderHandler *OrderHandlerImpl) ShipOrder(ctx *fiber.Ctx) error {
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	if err := orderHandler.OrderUseCase.ShipOrder(int64(orderID)); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "Order successfully marked as shipped.",
	})
}

func (orderHandler *OrderHandlerImpl) CheckOut(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := orderHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	if err := orderHandler.OrderUseCase.CheckOut(&body, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success checkout",
	})
}
func (orderHandler *OrderHandlerImpl) CheckOutConfirm(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqConfirmCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := orderHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := orderHandler.OrderUseCase.CheckOutConfirm(&body, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success confirm checkout",
		Data:    result,
	})
}
func (orderHandler *OrderHandlerImpl) GetOrderDetails(ctx *fiber.Ctx) error {
	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	order, err := orderHandler.OrderUseCase.GetOrderDetails(user.ID, int64(orderID))
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get all order",
		Data: map[string]any{
			"orders": order,
		},
	})
}

func (orderHandler *OrderHandlerImpl) GetAllOrder(ctx *fiber.Ctx) error {

	orders, err := orderHandler.OrderUseCase.GetAllOrder()
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get all order",
		Data: map[string]any{
			"orders": orders,
		},
	})
}
func (orderHandler *OrderHandlerImpl) GetAllUserOrder(ctx *fiber.Ctx) error {
	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}
	orders, err := orderHandler.OrderUseCase.GetAllUserOrder(user.ID)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get all order",
		Data: map[string]any{
			"orders": orders,
		},
	})
}

func (orderHandler *OrderHandlerImpl) GetLastDraftCheckOut(ctx *fiber.Ctx) error {

	user, err := orderHandler.getUserData(ctx)
	if err != nil {
		return err
	}

	result, err := orderHandler.OrderUseCase.GetLastDraftCheckOut(user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success get last draft checkout",
		Data: map[string]any{
			"checkout_draft": result,
		},
	})
}
