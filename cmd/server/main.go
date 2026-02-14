package main

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery"
	"ecommerce-system/internal/delivery/handler"
	"ecommerce-system/internal/delivery/middleware"
	"ecommerce-system/internal/infra/midtrans"
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
	config.InitMidtrans()

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

	userRepository := repository.NewUserRepository()
	productRepository := repository.NewProductRepository()
	categoryRepository := repository.NewCategoryRepository()
	addressRepository := repository.NewAddressRepository()
	cartRepository := repository.NewCartRepository()
	orderRepo := repository.NewOrderRepository()
	checkOutRepo := repository.NewCheckOutRepository()

	userUseCase := usecase.NewUserUseCase(userRepository, connection)
	addressUseCase := usecase.NewAddressUseCase(addressRepository, connection)
	productUseCase := usecase.NewProductUseCase(productRepository, connection)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepository, connection)
	cartUseCase := usecase.NewCartUseCase(cartRepository, productUseCase, connection)

	midtransPayment := midtrans.NewMidtransGateWay(config.Client)

	orderUseCase := usecase.NewOrderUseCase(checkOutRepo, orderRepo, cartRepository, productRepository, addressRepository, midtransPayment, connection)

	authHandler := handler.NewAuthController(userUseCase, validate, config)
	addressHandler := handler.NewAddressHandler(addressUseCase, validate)
	productHandler := handler.NewProductHandler(productUseCase, validate)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase, validate)
	cartHandler := handler.NewCartItemHandler(cartUseCase, validate)
	orderHandler := handler.NewOrderHandler(orderUseCase, config.Midtrans, validate)

	ctrl := &delivery.Handlers{
		Config:          config,
		AuthHandler:     authHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
		AddressHandler:  addressHandler,
		CartHandler:     cartHandler,
		OrderHandler:    orderHandler,
	}
	delivery.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
