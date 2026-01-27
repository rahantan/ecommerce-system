package routes

import (
	"ecommerce-system/config"
	addresshandlers "ecommerce-system/internal/handlers/addresses"
	authhandlers "ecommerce-system/internal/handlers/auth"
	categoryhandlers "ecommerce-system/internal/handlers/categories"
	producthandlers "ecommerce-system/internal/handlers/products"
	"ecommerce-system/internal/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Handlers struct {
	*config.Config
	authhandlers.AuthHandlers
	categoryhandlers.CategoryHandlers
	producthandlers.ProductHandlers
	addresshandlers.AddressHandlers
}

func NewRoute(app *fiber.App, handler *Handlers) {

	limit := limiter.New(middlewares.Limiter())

	// AUTH (guest)
	auth := app.Group("/auth", limit)
	auth.Post("/login", handler.AuthHandlers.Login)
	auth.Post("/register", handler.AuthHandlers.Register)

	// PUBLIC API
	public := app.Group("/api", limit)
	public.Get("/products/:productId", handler.ProductHandlers.GetProductById)
	public.Get("/products", handler.ProductHandlers.GetAllProduct)

	public.Get("/categories/:categoryId", handler.CategoryHandlers.GetCategoryById)
	public.Get("/categories", handler.CategoryHandlers.GetAllCategory)

	// PRIVATE API
	private := app.Group(
		"/api",
		middlewares.JwtValidationToken(handler.Config.Jwt.SecretKey),
		limit,
	)
	private.Delete("/logout", handler.AuthHandlers.Logout)

	// ADMIN
	admin := private.Group("/admin", middlewares.Authorization(1))
	admin.Post("/products", handler.ProductHandlers.CreateProduct)
	admin.Put("/products", handler.ProductHandlers.UpdateProductById)

	admin.Post("/categories", handler.CategoryHandlers.CreateCategory)
	admin.Put("/categories", handler.CategoryHandlers.UpdateCategoryById)

	// CUSTOMER
	customer := private.Group("/customers", middlewares.Authorization(2))
	customer.Get("/addresses/active", handler.AddressHandlers.GetUserActiveAddress)
	customer.Get("/addresses", handler.AddressHandlers.GetAllAddress)
	customer.Post("/addresses", handler.AddressHandlers.CreateAddress)
	customer.Put("/addresses/:addressId", handler.AddressHandlers.UpdateAddressByUserId)
}
