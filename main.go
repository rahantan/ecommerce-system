package main

import (
	"ecommerce-system/config"
	"reflect"
	"strings"

	"ecommerce-system/internal/exceptions"
	authhandlers "ecommerce-system/internal/handlers/auth"

	userrepositories "ecommerce-system/internal/repositories/users"
	"ecommerce-system/internal/routes"
	authservices "ecommerce-system/internal/services/auth"
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
	userService := userservices.NewUserService(userRepository)

	authService := authservices.NewAuthService(userService)

	authHandler := authhandlers.NewAuthController(authService, validate, config)

	ctrl := &routes.Handlers{
		Config:       config,
		AuthHandlers: authHandler,
	}
	routes.NewRoute(app, ctrl)
	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
