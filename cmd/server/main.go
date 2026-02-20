package main

import (
	"ecommerce-system/config"
	"ecommerce-system/internal/delivery"
	cartHandler "ecommerce-system/internal/delivery/handler/cart"
	categoryHandler "ecommerce-system/internal/delivery/handler/category"
	coHandler "ecommerce-system/internal/delivery/handler/checkout"
	orderHandler "ecommerce-system/internal/delivery/handler/order"
	productHandler "ecommerce-system/internal/delivery/handler/product"
	userHandler "ecommerce-system/internal/delivery/handler/user"

	"ecommerce-system/internal/delivery/middleware"
	"ecommerce-system/internal/infra/external/midtrans"
	"ecommerce-system/internal/infra/persistence"

	cartUsecase "ecommerce-system/internal/usecase/cart"
	categoryUsecase "ecommerce-system/internal/usecase/category"
	coUsecase "ecommerce-system/internal/usecase/checkout"
	orderUsecase "ecommerce-system/internal/usecase/order"
	productUsecase "ecommerce-system/internal/usecase/product"
	userUsecase "ecommerce-system/internal/usecase/user"
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

	userRepository := persistence.NewUserRepository()
	productRepository := persistence.NewProductRepository()
	categoryRepository := persistence.NewCategoryRepository()
	addressRepository := persistence.NewAddressRepository()
	cartRepository := persistence.NewCartRepository()
	orderRepo := persistence.NewOrderRepository()
	checkOutRepo := persistence.NewCheckOutRepository()
	paymentRepo := persistence.NewPaymentRepository()

	userUseCase := userUsecase.NewUserUseCase(userRepository, addressRepository, connection)

	// addressUseCase := usecase.NewAddressUseCase(addressRepository, connection)
	productUseCase := productUsecase.NewProductUseCase(productRepository, connection)
	categoryUseCase := categoryUsecase.NewCategoryUseCase(categoryRepository, connection)
	cartUseCase := cartUsecase.NewCartUseCase(cartRepository, productUseCase, connection)

	midtransPayment := midtrans.NewMidtransGateWay(config.Client)
	checkOutUseCase := coUsecase.NewCheckOutUseCase(checkOutRepo, orderRepo, cartRepository, productRepository, addressRepository, paymentRepo, midtransPayment, connection)

	orderUseCase := orderUsecase.NewOrderUseCase(checkOutRepo, orderRepo, cartRepository, productRepository, addressRepository, paymentRepo, midtransPayment, connection)

	userHandler := userHandler.NewUserHandler(userUseCase, validate, config)
	productHandler := productHandler.NewProductHandler(productUseCase, validate)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryUseCase, validate)
	cartHandler := cartHandler.NewCartItemHandler(cartUseCase, validate)
	orderHandler := orderHandler.NewOrderHandler(orderUseCase, config.Midtrans, validate)
	checkOutHandler := coHandler.NewCheckOutHandler(checkOutUseCase, validate)

	ctrl := &delivery.Handlers{
		Config:          config,
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  productHandler,
		CartHandler:     cartHandler,
		OrderHandler:    orderHandler,
		CheckOutHandler: checkOutHandler,
	}
	delivery.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
