package order

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/pkg"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (orderHandler *OrderHandlerImpl) WebHookMidtransNotif(ctx *fiber.Ctx) error {
	var body request.MidtransNotification
	if err := ctx.BodyParser(&body); err != nil {
		fmt.Println(err.Error())
		return ctx.SendStatus(pkg.ErrCustomInvalidPayload.GetStatusCode())
	}
	fmt.Println("MASUK WEB HOOK")
	fmt.Println("Order ID: ", body.OrderID)
	fmt.Println("STATUS PAYMENT: ", body.TransactionStatus)
	// return ctx.SendStatus(200)

	if !orderHandler.midTransVerify(body.OrderID, body.StatusCode, body.GrossAmount, body.SignatureKey) {
		fmt.Println("invalid midtrans verify")
		return ctx.SendStatus(pkg.ErrCustomForbidden.GetStatusCode())
	}

	if body.FraudStatus != "accept" {
		fmt.Println("ente mau nipu yah?")
		return ctx.SendStatus(pkg.ErrCustomForbidden.GetStatusCode())
	}

	orderID, err := strconv.Atoi(body.OrderID)
	if err != nil {
		fmt.Println(err.Error())
		return pkg.ErrCustomInvalidOrderId
	}

	if err := orderHandler.OrderUseCase.UpdateStatusPayment(int64(orderID), body.TransactionStatus); err != nil {
		return err
	}

	if err := orderHandler.UpdateStatusPayment(int64(orderID), body.TransactionStatus); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
