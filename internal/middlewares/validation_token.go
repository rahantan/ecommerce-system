package middlewares

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtValidationToken(secretKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenStr := ctx.Cookies("token")
		if tokenStr == "" {
			return exceptions.ErrCustomTokenEmpty
		}

		claims := &utils.JwtPayload{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signature algorithm")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			return err
		}
		if !token.Valid {
			return errors.New("invalid parse token")
		}
		ctx.Locals("user", claims.ResUser)
		return ctx.Next()
	}
}
