package routes

import (
	"ecommerce-system/config"
	authhandlers "ecommerce-system/internal/handlers/auth"
	"ecommerce-system/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	*config.Config
	authhandlers.AuthHandlers
}

func NewRoute(app *fiber.App, ctrl *Handlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	auth := app.Group("/auth")
	auth.Post("/login", ctrl.Login)
	auth.Post("/register", ctrl.Register)

	private := app.Group("/api", middlewares.JwtValidationToken(ctrl.Config.Jwt.SecretKey))

	private.Get("/home", func(c *fiber.Ctx) error {
		return c.SendString("hello private")
	})

	private.Delete("/logout", ctrl.Logout)
}
