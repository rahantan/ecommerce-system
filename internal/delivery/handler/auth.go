package handler

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/dto/response"
	"ecommerce-system/internal/domain"
	"ecommerce-system/internal/infra/jwt"
	"ecommerce-system/internal/pkg"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerImpl struct {
	// authservices.AuthServices
	domain.UserUseCase
	*validator.Validate
	*config.Config
}

func NewAuthController(user domain.UserUseCase, v *validator.Validate, c *config.Config) domain.AuthHandler {
	return &AuthHandlerImpl{
		UserUseCase: user,
		Validate:    v,
		Config:      c,
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

	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := auth.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := auth.UserUseCase.Login(&body)
	if err != nil {
		return err
	}

	token, err := jwt.GetToken(*result, auth.Config.Jwt.SecretKey)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(24 * time.Hour),
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
	if err := ctx.BodyParser(&body); err != nil {
		return pkg.ErrCustomInvalidPayload
	}

	if err := auth.Validate.Struct(&body); err != nil {
		return pkg.ValidationError(err)
	}

	result, err := auth.UserUseCase.Register(&body)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.ResponseStandard{
		Success: true,
		Message: "success register",
		Data: map[string]any{
			"user": result,
		},
	})
}
