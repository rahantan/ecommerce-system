package main

import (
	"ecommerce-system/config"
	"reflect"
	"strings"

	"ecommerce-system/internal/exceptions"
	addresshandlers "ecommerce-system/internal/handlers/addresses"
	authhandlers "ecommerce-system/internal/handlers/auth"
	categoryhandlers "ecommerce-system/internal/handlers/categories"
	producthandlers "ecommerce-system/internal/handlers/products"

	addressrepositories "ecommerce-system/internal/repositories/addresses"
	categoryrepositories "ecommerce-system/internal/repositories/categories"
	productrepositories "ecommerce-system/internal/repositories/products"
	userrepositories "ecommerce-system/internal/repositories/users"
	"ecommerce-system/internal/routes"
	addressservices "ecommerce-system/internal/services/addresses"
	authservices "ecommerce-system/internal/services/auth"
	categoryservices "ecommerce-system/internal/services/categories"
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

	userService := userservices.NewUserService(userRepository)
	addressService := addressservices.NewAddressService(addressRepository)
	authService := authservices.NewAuthService(userService)
	productService := productservices.NewProductService(productRepository)
	categoryService := categoryservices.NewCategoryService(categoryRepository)

	authHandler := authhandlers.NewAuthController(authService, validate, config)
	addressHandler := addresshandlers.NewAddressHandler(addressService, validate)
	productHandler := producthandlers.NewProductHandler(productService, validate)
	categoryHandler := categoryhandlers.NewCategoryHandler(categoryService, validate)

	ctrl := &routes.Handlers{
		Config:           config,
		AuthHandlers:     authHandler,
		CategoryHandlers: categoryHandler,
		ProductHandlers:  productHandler,
		AddressHandlers:  addressHandler,
	}
	routes.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
