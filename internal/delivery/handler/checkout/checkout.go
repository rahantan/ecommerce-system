package checkout

import (
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CheckOutHandlerImpl struct {
	domain.CheckOutUseCase
	*validator.Validate
}

func NewCheckOutHandler(order domain.CheckOutUseCase, v *validator.Validate) domain.CheckOutHandler {
	return &CheckOutHandlerImpl{
		CheckOutUseCase: order,
		Validate:        v,
	}
}

func (coHandler *CheckOutHandlerImpl) CheckOut(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqCheckout
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := coHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	if err := coHandler.CheckOutUseCase.CheckOut(&body, user.ID); err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusCreated, "success checkout", nil)

}

func (coHandler *CheckOutHandlerImpl) CheckOutConfirm(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	var body request.ReqConfirmCheckout

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := coHandler.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := coHandler.CheckOutUseCase.CheckOutConfirm(&body, user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusCreated, "success confirm checkout", map[string]any{
		"payments": result,
	})
}

func (coHandler *CheckOutHandlerImpl) GetLastDraftCheckOut(ctx *fiber.Ctx) error {

	user, err := handler.GetUserData(ctx)
	if err != nil {
		return err
	}

	result, err := coHandler.CheckOutUseCase.GetLastDraftCheckOut(user.ID)
	if err != nil {
		return err
	}

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success get last draft checkout", map[string]any{
		"checkouts": result,
	})
}
