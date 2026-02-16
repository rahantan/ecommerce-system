package delivery

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery/middleware"
	"ecommerce-system/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Handlers struct {
	*config.Config
	domain.AuthHandler
	domain.CategoryHandler
	domain.ProductHandler
	domain.AddressHandler
	domain.CartHandler
	domain.OrderHandler
	domain.MidtransGateWay
}

func NewRoute(app *fiber.App, handler *Handlers) {

	limit := limiter.New(middleware.Limiter())

	api := app.Group("/api")
	webHook := api.Group("/webhook")
	handler.WebHook(webHook)

	public := api.Group("/public", limit)
	handler.Public(public)

	private := api.Group("/private", middleware.JwtValidationToken(handler.Config.Jwt.SecretKey), limit)
	private.Delete("/logout", handler.AuthHandler.Logout)

	admin := private.Group("/admin", middleware.Authorization(1))
	handler.Admin(admin)

	customers := private.Group("/customers", middleware.Authorization(2))
	handler.Customers(customers)

}
func (route *Handlers) WebHook(webHook fiber.Router) {
	webHook.Post("/midtrans/notif", route.OrderHandler.WebHookMidtransNotif)
}
func (route *Handlers) Public(public fiber.Router) {
	auth := public.Group("/auth")
	auth.Post("/login", route.AuthHandler.Login)
	auth.Post("/register", route.AuthHandler.Register)

	public.Get("/products/:productId", route.ProductHandler.GetProductById)
	public.Get("/products", route.ProductHandler.GetAllProduct)

	public.Get("/categories/:categoryId", route.CategoryHandler.GetCategoryById)
	public.Get("/categories", route.CategoryHandler.GetAllCategory)
}

func (route *Handlers) Admin(admin fiber.Router) {
	admin.Put("/products/:productId", route.ProductHandler.UpdateProductById)
	admin.Post("/products", route.ProductHandler.CreateProduct)

	admin.Put("/categories/:categoryId", route.CategoryHandler.UpdateCategoryById)
	admin.Post("/categories", route.CategoryHandler.CreateCategory)
}

func (route *Handlers) Customers(customer fiber.Router) {
	customer.Get("/addresses/active", route.AddressHandler.GetUserActiveAddress)
	customer.Get("/addresses", route.AddressHandler.GetAllAddress)
	customer.Post("/addresses", route.AddressHandler.CreateAddress)
	customer.Put("/addresses/:addressId", route.AddressHandler.UpdateAddressByUserId)

	customer.Post("/carts", route.CartHandler.AddCartItem)
	customer.Get("/carts", route.CartHandler.GetAllUserCartItem)
	customer.Delete("/carts", route.CartHandler.DeleteCartItemsByIDs)

	customer.Get("/orders/checkout", route.OrderHandler.GetLastDraftCheckOut)
	customer.Post("/orders/checkout", route.OrderHandler.CheckOut)
	customer.Post("/orders/confirm", route.OrderHandler.CheckOutConfirm)

	customer.Get("/orders/:orderID", route.OrderHandler.GetOrderDetails)
	customer.Get("/orders", route.OrderHandler.GetAllOrder)
}
