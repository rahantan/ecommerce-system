package middleware

import (
	jwtauth "ecommerce-system/internal/infra/jwt"
	"ecommerce-system/internal/pkg"
	"errors"

	"github.com/gofiber/fiber/v2"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

func JwtValidationToken(secretKey string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenStr := ctx.Cookies("token")
		if tokenStr == "" {
			return pkg.NewError(pkg.KindUnauthorized, "token is required", nil)
		}

		claims := &jwtauth.JwtPayload{}
		token, err := jwtlib.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (any, error) {
			if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signature algorithm")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			return pkg.NewError(pkg.KindUnauthorized, err.Error(), nil)
		}

		if !token.Valid {
			return pkg.NewError(pkg.KindUnauthorized, "invalid parse token", nil)
		}

		ctx.Locals("user", claims.ResUser)

		return ctx.Next()
	}
}
