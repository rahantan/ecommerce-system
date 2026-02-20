package order

import (
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/pkg"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (orderHandler *OrderHandlerImpl) GetUserPaymentByOrderID(ctx *fiber.Ctx) error {
	user, err := handler.GetUserData(ctx)
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

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get payment", map[string]any{
		"payments": result,
	})

}
