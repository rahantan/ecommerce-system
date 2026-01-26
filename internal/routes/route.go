package routes

import (
	"ecommerce-system/config"
	authhandlers "ecommerce-system/internal/handlers/auth"
	categoryhandlers "ecommerce-system/internal/handlers/categories"
	producthandlers "ecommerce-system/internal/handlers/products"
	"ecommerce-system/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	*config.Config
	authhandlers.AuthHandlers
	categoryhandlers.CategoryHandlers
	producthandlers.ProductHandlers
}

func NewRoute(app *fiber.App, hndlr *Handlers) {

	auth := app.Group("/auth")

	auth.Post("/login", hndlr.AuthHandlers.Login)
	auth.Post("/register", hndlr.AuthHandlers.Register)

	public := app.Group("/api")

	public.Get("/products/:productId", hndlr.ProductHandlers.GetProductById)
	public.Get("/products", hndlr.ProductHandlers.GetAllProduct)

	public.Get("/categories/:categoryId", hndlr.GetCategoryById)
	public.Get("/categories", hndlr.GetAllCategory)

	private := app.Group("/api", middlewares.JwtValidationToken(hndlr.Config.Jwt.SecretKey))
	private.Delete("/logout", hndlr.AuthHandlers.Logout)
	//ADMIN ENDPOINT
	admin := private.Group("", middlewares.Authorization(1))
	admin.Post("/products", hndlr.ProductHandlers.CreateProduct)
	admin.Put("/products", hndlr.ProductHandlers.UpdateProductById)

	admin.Post("/categories", hndlr.CategoryHandlers.CreateCategory)
	admin.Put("/categories", hndlr.CategoryHandlers.UpdateCategoryById)

	//CUSTOMER ENDPOINT
	// customer := private.Group("", middlewares.Authorization(2))

}
