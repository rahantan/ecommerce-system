package authhandlers

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	authservices "ecommerce-system/internal/services/auth"
	"ecommerce-system/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerImpl struct {
	authservices.AuthServices
	*validator.Validate
	*config.Config
}

func NewAuthController(auth authservices.AuthServices, v *validator.Validate, c *config.Config) AuthHandlers {
	return &AuthHandlerImpl{
		AuthServices: auth,
		Validate:     v,
		Config:       c,
	}
}

func (auth *AuthHandlerImpl) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success log out",
	})
}
func (auth *AuthHandlerImpl) Login(ctx *fiber.Ctx) error {
	var body request.ReqLogin

	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.UpdateMessageErr(exceptions.ErrCustomInvalidPayload, exceptions.MsgFailLogin)
	}

	err = auth.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailLogin)
	}

	result, err := auth.AuthServices.Login(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailLogin)
	}

	token, err := utils.GetToken(*result, auth.Config.Jwt.SecretKey)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailLogin)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})

	return ctx.Status(fiber.StatusOK).JSON(response.ResponseStandard{
		Success: true,
		Message: "success login",
		Data: map[string]any{
			"user": result,
		},
	})
}
func (auth *AuthHandlerImpl) Register(ctx *fiber.Ctx) error {
	var body request.ReqCreateUser

	err := ctx.BodyParser(&body)
	if err != nil {
		return utils.UpdateMessageErr(exceptions.ErrCustomInvalidPayload, exceptions.MsgFailRegister)
	}

	err = auth.Validate.Struct(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailRegister)
	}

	result, err := auth.AuthServices.Register(&body)
	if err != nil {
		return utils.UpdateMessageErr(err, exceptions.MsgFailRegister)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success register",
		Data: map[string]any{
			"user": result,
		},
	})
}
