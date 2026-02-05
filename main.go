package main

import (
	"ecommerce-system/config"
	"reflect"
	"strings"

	"ecommerce-system/internal/exceptions"
	addresshandlers "ecommerce-system/internal/handlers/addresses"
	authhandlers "ecommerce-system/internal/handlers/auth"
	carthandlers "ecommerce-system/internal/handlers/carts"
	categoryhandlers "ecommerce-system/internal/handlers/categories"
	orderhandlers "ecommerce-system/internal/handlers/orders"
	producthandlers "ecommerce-system/internal/handlers/products"

	addressrepositories "ecommerce-system/internal/repositories/addresses"
	cartitemrepositories "ecommerce-system/internal/repositories/carts"
	categoryrepositories "ecommerce-system/internal/repositories/categories"
	orderrepositories "ecommerce-system/internal/repositories/orders"
	productrepositories "ecommerce-system/internal/repositories/products"
	userrepositories "ecommerce-system/internal/repositories/users"
	"ecommerce-system/internal/routes"
	addressservices "ecommerce-system/internal/services/addresses"
	cartitemsservices "ecommerce-system/internal/services/carts"
	categoryservices "ecommerce-system/internal/services/categories"
	orderservices "ecommerce-system/internal/services/orders"
	productservices "ecommerce-system/internal/services/products"
	userservices "ecommerce-system/internal/services/users"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config := config.LoadConfig()
	connection := config.ConnectionDb()
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Tag.Get("json")
		if tag == "-" {
			return ""
		}
		return strings.Split(tag, ",")[0]
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: exceptions.ErrorHandler,
	})

	userRepository := userrepositories.NewUserRepository(connection)
	productRepository := productrepositories.NewProductRepository(connection)
	categoryRepository := categoryrepositories.NewCategoryRepository(connection)
	addressRepository := addressrepositories.NewAddressRepository(connection)
	cartRepository := cartitemrepositories.NewCartItemRepository()
	orderRepo := orderrepositories.NewOrderRepository()
	checkOutRepo := orderrepositories.NewCheckoutSession()

	userService := userservices.NewUserService(userRepository, connection)
	addressService := addressservices.NewAddressService(addressRepository, connection)
	productService := productservices.NewProductService(productRepository, connection)
	categoryService := categoryservices.NewCategoryService(categoryRepository, connection)
	cartService := cartitemsservices.NewCartItemService(cartRepository, productService, connection)

	orderService := orderservices.NewOrderService(checkOutRepo, orderRepo, cartRepository, productRepository, addressRepository, connection)

	authHandler := authhandlers.NewAuthController(userService, validate, config)
	addressHandler := addresshandlers.NewAddressHandler(addressService, validate)
	productHandler := producthandlers.NewProductHandler(productService, validate)
	categoryHandler := categoryhandlers.NewCategoryHandler(categoryService, validate)
	cartHandler := carthandlers.NewCartItemHandler(cartService, validate)
	orderHandler := orderhandlers.NewOrderHandler(orderService, validate)

	ctrl := &routes.Handlers{
		Config:           config,
		AuthHandlers:     authHandler,
		CategoryHandlers: categoryHandler,
		ProductHandlers:  productHandler,
		AddressHandlers:  addressHandler,
		CartItemHandlers: cartHandler,
		OrderHandlers:    orderHandler,
	}
	routes.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
