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

	public := api.Group("/public", limit)
	handler.Public(public)

	private := api.Group("/private", middlewares.JwtValidationToken(handler.Config.Jwt.SecretKey), limit)
	private.Delete("/logout", handler.AuthHandlers.Logout)

	admin := private.Group("/admin", middlewares.Authorization(1))
	handler.Admin(admin)

	customers := private.Group("/customers", middlewares.Authorization(2))
	handler.Customers(customers)

}

func (route *Handlers) Public(public fiber.Router) {
	auth := public.Group("/auth")
	auth.Post("/login", route.AuthHandlers.Login)
	auth.Post("/register", route.AuthHandlers.Register)

	public.Get("/products/:productId", route.ProductHandlers.GetProductById)
	public.Get("/products", route.ProductHandlers.GetAllProduct)

	public.Get("/categories/:categoryId", route.CategoryHandlers.GetCategoryById)
	public.Get("/categories", route.CategoryHandlers.GetAllCategory)
}

func (route *Handlers) Admin(admin fiber.Router) {
	admin.Put("/products/:productId", route.ProductHandlers.UpdateProductById)
	admin.Post("/products", route.ProductHandlers.CreateProduct)

	admin.Put("/categories/:categoryId", route.CategoryHandlers.UpdateCategoryById)
	admin.Post("/categories", route.CategoryHandlers.CreateCategory)

}

func (route *Handlers) Customers(customer fiber.Router) {
	customer.Get("/addresses/active", route.AddressHandlers.GetUserActiveAddress)
	customer.Get("/addresses", route.AddressHandlers.GetAllAddress)
	customer.Post("/addresses", route.AddressHandlers.CreateAddress)
	customer.Put("/addresses/:addressId", route.AddressHandlers.UpdateAddressByUserId)

	customer.Post("/carts", route.CartItemHandlers.AddCartItem)
	customer.Get("/carts", route.CartItemHandlers.GetAllUserCartItem)
	customer.Delete("/carts", route.CartItemHandlers.DeleteCartItemsByIDs)

	customer.Get("/order/checkout", route.OrderHandlers.GetLastDraftCheckOut)
	customer.Post("/order/checkout", route.OrderHandlers.CheckOut)
	customer.Post("/order/confirm", route.OrderHandlers.CheckOutConfirm)
}
