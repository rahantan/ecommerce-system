package order

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"
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

func (orderHandler *OrderHandlerImpl) ReceiveOrder(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
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

	return handler.SuccessResponse(ctx, fiber.StatusOK, "Order successfully marked as received.", nil)

}

func (orderHandler *OrderHandlerImpl) ShipOrder(ctx *fiber.Ctx) error {

	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	if err := orderHandler.OrderUseCase.ShipOrder(int64(orderID)); err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "Order successfully marked as shipped.", nil)

}

func (orderHandler *OrderHandlerImpl) GetOrderDetails(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}
	orderID, err := strconv.Atoi(ctx.Params("orderID"))
	if err != nil {
		return pkg.ErrCustomInvalidOrderId
	}

	result, err := orderHandler.OrderUseCase.GetOrderDetails(user.ID, int64(orderID))
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get order details", map[string]any{
		"orders": result,
	})

}

func (orderHandler *OrderHandlerImpl) GetAllOrder(ctx *fiber.Ctx) error {

	result, err := orderHandler.OrderUseCase.GetAllOrder()
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get orders", map[string]any{
		"orders": result,
	})

}

func (orderHandler *OrderHandlerImpl) GetUserOrders(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	result, err := orderHandler.OrderUseCase.GetAllOrderByUserID(user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get user orders", map[string]any{
		"orders": result,
	})

}
