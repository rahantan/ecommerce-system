package utils

import (
	"ecommerce-system/internal/dto/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	jwt.RegisteredClaims
	response.ResUser
}

func GetToken(user response.ResUser, secretKey string) (string, error) {
	start := time.Now()
	payload := JwtPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(start.Add(time.Duration(time.Hour * 12))),
		},
		ResUser: user,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
