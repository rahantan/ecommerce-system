package routes

import (
	"ecommerce-system/config"
	addresshandlers "ecommerce-system/internal/handlers/addresses"
	authhandlers "ecommerce-system/internal/handlers/auth"
	carthandlers "ecommerce-system/internal/handlers/carts"
	categoryhandlers "ecommerce-system/internal/handlers/categories"
	orderhandlers "ecommerce-system/internal/handlers/orders"
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
	carthandlers.CartItemHandlers
	orderhandlers.OrderHandlers
}

func NewRoute(app *fiber.App, handler *Handlers) {

	limit := limiter.New(middlewares.Limiter())

	api := app.Group("/api")

	// PUBLIC API
	public := api.Group("/public", limit)
	// AUTH
	auth := public.Group("/auth")
	auth.Post("/login", handler.AuthHandlers.Login)
	auth.Post("/register", handler.AuthHandlers.Register)

	// public := app.Group("/api", limit)
	public.Get("/products/:productId", handler.ProductHandlers.GetProductById)
	public.Get("/products", handler.ProductHandlers.GetAllProduct)

	public.Get("/categories/:categoryId", handler.CategoryHandlers.GetCategoryById)
	public.Get("/categories", handler.CategoryHandlers.GetAllCategory)

	// PRIVATE API
	private := api.Group("/private", middlewares.JwtValidationToken(handler.Config.Jwt.SecretKey), limit)
	private.Delete("/logout", handler.AuthHandlers.Logout)

	// ADMIN
	admin := private.Group("/admin", middlewares.Authorization(1))
	admin.Put("/products/:productId", handler.ProductHandlers.UpdateProductById)
	admin.Post("/products", handler.ProductHandlers.CreateProduct)

	admin.Put("/categories/:categoryId", handler.CategoryHandlers.UpdateCategoryById)
	admin.Post("/categories", handler.CategoryHandlers.CreateCategory)

	// CUSTOMER
	customer := private.Group("/customers", middlewares.Authorization(2))

	customer.Get("/addresses/active", handler.AddressHandlers.GetUserActiveAddress)
	customer.Get("/addresses", handler.AddressHandlers.GetAllAddress)
	customer.Post("/addresses", handler.AddressHandlers.CreateAddress)
	customer.Put("/addresses/:addressId", handler.AddressHandlers.UpdateAddressByUserId)

	customer.Post("/carts", handler.CartItemHandlers.AddCartItem)
	customer.Get("/carts", handler.CartItemHandlers.GetAllUserCartItem)
	customer.Delete("/carts", handler.CartItemHandlers.DeleteCartItemsByIDs)

	customer.Get("/order/checkout", handler.OrderHandlers.GetLastDraftCheckOut)
	customer.Post("/order/checkout", handler.OrderHandlers.CheckOut)
	customer.Post("/order/confirm", handler.OrderHandlers.CheckOutConfirm)
}
