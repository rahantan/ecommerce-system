package authhandlers

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/dto/request"
	"ecommerce-system/internal/dto/response"
	"ecommerce-system/internal/exceptions"
	userservices "ecommerce-system/internal/services/users"
	"ecommerce-system/internal/utils"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerImpl struct {
	// authservices.AuthServices
	userservices.UserServices
	*validator.Validate
	*config.Config
}

func NewAuthController(user userservices.UserServices, v *validator.Validate, c *config.Config) AuthHandlers {
	return &AuthHandlerImpl{
		UserServices: user,
		Validate:     v,
		Config:       c,
	}
}

func (auth *AuthHandlerImpl) withMessage(err error, msg string) error {
	return utils.WithMessage(err, msg)
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
		return auth.withMessage(exceptions.ErrCustomInvalidPayload, "failed to login")
	}

	err = auth.Validate.Struct(&body)
	if err != nil {
		return auth.withMessage(exceptions.ValidationError(err), "failed to login")
	}

	result, err := auth.UserServices.Login(&body)
	if err != nil {
		return auth.withMessage(err, "failed to login")
	}

	token, err := utils.GetToken(*result, auth.Config.Jwt.SecretKey)
	if err != nil {
		return auth.withMessage(err, "failed to login")
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
		return auth.withMessage(exceptions.ErrCustomInvalidPayload, "failed to register")
	}

	err = auth.Validate.Struct(&body)
	if err != nil {
		return auth.withMessage(exceptions.ValidationError(err), "failed to register")
	}

	result, err := auth.UserServices.Register(&body)
	if err != nil {
		return auth.withMessage(err, "failed to register")
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success register",
		Data: map[string]any{
			"user": result,
		},
	})
}
