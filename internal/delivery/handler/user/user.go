package user

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery/dto/request"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/domain"
	jwtAuth "ecommerce-system/internal/infra/external/jwt"
	"ecommerce-system/internal/pkg"

	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandlerImpl struct {
	domain.UserUseCase
	*validator.Validate
	*config.Config
}

func NewUserHandler(user domain.UserUseCase, v *validator.Validate, c *config.Config) domain.UserHandler {
	return &UserHandlerImpl{
		UserUseCase: user,
		Validate:    v,
		Config:      c,
	}
}

func (auth *UserHandlerImpl) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:   "token",
		MaxAge: -1,
	})

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success log out", nil)
}

func (auth *UserHandlerImpl) Login(ctx *fiber.Ctx) error {

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

	token, err := jwtAuth.GetToken(*result, auth.Config.Jwt.SecretKey)
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

	return handler.SuccessResponse(ctx, fiber.StatusOK, "success login", map[string]any{
		"users": result,
	})

}

func (auth *UserHandlerImpl) Register(ctx *fiber.Ctx) error {

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

	return handler.SuccessResponse(ctx, fiber.StatusCreated, "success register", map[string]any{
		"users": result,
	})

}
