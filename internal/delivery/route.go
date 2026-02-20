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
	domain.UserHandler
	domain.CategoryHandler
	domain.ProductHandler
	domain.CartHandler
	domain.OrderHandler
	domain.MidtransGateWay
	domain.CheckOutHandler
}

func NewRoute(app *fiber.App, handler *Handlers) {

	limit := limiter.New(middleware.Limiter())

	api := app.Group("/api")
	webHook := api.Group("/webhook")
	handler.WebHook(webHook)

	public := api.Group("/public", limit)
	handler.Public(public)

	private := api.Group("/private", middleware.JwtValidationToken(handler.Config.Jwt.SecretKey), limit)
	private.Delete("/logout", handler.UserHandler.Logout)

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
	auth.Post("/login", route.UserHandler.Login)
	auth.Post("/register", route.UserHandler.Register)

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

	admin.Put("/orders/:orderID/ship", route.OrderHandler.ShipOrder)
	admin.Get("/orders", route.OrderHandler.GetAllOrder)
}

func (route *Handlers) Customers(customer fiber.Router) {
	customer.Get("/addresses/active", route.UserHandler.GetUserActiveAddress)
	customer.Get("/addresses", route.UserHandler.GetAllAddress)
	customer.Post("/addresses", route.UserHandler.CreateAddress)
	customer.Put("/addresses/:addressId", route.UserHandler.UpdateAddressByUserId)

	customer.Put("/carts/:cartID", route.CartHandler.UpdateCartItemByID)
	customer.Post("/carts", route.CartHandler.AddCartItem)
	customer.Get("/carts", route.CartHandler.GetAllUserCartItem)
	customer.Delete("/carts", route.CartHandler.DeleteCartItemsByIDs)

	customer.Get("/orders/checkout", route.CheckOutHandler.GetLastDraftCheckOut)
	customer.Post("/orders/checkout", route.CheckOutHandler.CheckOut)
	customer.Post("/orders/confirm", route.CheckOutHandler.CheckOutConfirm)

	customer.Put("/orders/:orderID/receive", route.OrderHandler.ReceiveOrder)
	customer.Get("/orders/:orderID/payment", route.OrderHandler.GetUserPaymentByOrderID)
	customer.Get("/orders/:orderID", route.OrderHandler.GetOrderDetails)
	customer.Get("/orders", route.OrderHandler.GetUserOrders)
}
