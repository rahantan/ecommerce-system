package middlewares

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JwtValidationToken(secretKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenStr := ctx.Cookies("token")
		if tokenStr == "" {
			return exceptions.NewError(exceptions.DefaultMsgUnauthorized, "token is required", http.StatusUnauthorized)
		}

		claims := &utils.JwtPayload{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signature algorithm")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			return exceptions.NewError(exceptions.DefaultMsgUnauthorized, "invalid token or expired", http.StatusUnauthorized)
		}
		if !token.Valid {
			return errors.New("invalid parse token")
		}
		ctx.Locals("user", claims.ResUser)
		return ctx.Next()
	}
}
