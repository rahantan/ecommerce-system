package middleware

import (
	"ecommerce-system/internal/delivery/dto/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func Limiter() limiter.Config {
	return limiter.Config{
		Max:        5,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			userId, ok := c.Locals("user").(response.ResUser)
			if ok {
				return strconv.Itoa(int(userId.ID))
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(response.ResponseStandard{
				Success: false,
				Message: "Too Many Requests",
			})
		},
	}
}
