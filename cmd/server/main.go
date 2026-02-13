package main

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/delivery/middleware"
	"ecommerce-system/internal/repository"
	"ecommerce-system/internal/usecase"
	"fmt"
	"reflect"
	"strings"

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
		ErrorHandler: middleware.ErrorHandler,
	})

	userRepository := repository.NewUserRepository(connection)
	productRepository := repository.NewProductRepository(connection)
	categoryRepository := repository.NewCategoryRepository(connection)
	addressRepository := repository.NewAddressRepository(connection)
	cartRepository := repository.NewCartItemRepository()
	orderRepo := repository.NewOrderRepository()
	checkOutRepo := repository.NewCheckoutSession()

	userService := usecase.NewUserService(userRepository, connection)
	addressService := usecase.NewAddressService(addressRepository, connection)
	productService := usecase.NewProductService(productRepository, connection)
	categoryService := usecase.NewCategoryService(categoryRepository, connection)
	cartService := usecase.NewCartItemService(cartRepository, productService, connection)

	orderService := usecase.NewOrderService(checkOutRepo, orderRepo, cartRepository, productRepository, addressRepository, connection)

	authHandler := handler.NewAuthController(userService, validate, config)
	addressHandler := handler.NewAddressHandler(addressService, validate)
	productHandler := handler.NewProductHandler(productService, validate)
	categoryHandler := handler.NewCategoryHandler(categoryService, validate)
	cartHandler := handler.NewCartItemHandler(cartService, validate)
	orderHandler := handler.NewOrderHandler(orderService, validate)

	ctrl := &delivery.Handlers{
		Config:           config,
		AuthHandlers:     authHandler,
		CategoryHandlers: categoryHandler,
		ProductHandlers:  productHandler,
		AddressHandlers:  addressHandler,
		CartItemHandlers: cartHandler,
		OrderHandlers:    orderHandler,
	}
	delivery.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
