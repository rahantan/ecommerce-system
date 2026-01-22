package main

import (
	"ecommerce-system/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config := config.LoadConfig()
	app := fiber.New()

	if err := app.Listen(fmt.Sprintf("%s:%s",
		config.Server.Host,
		config.Server.Port,
	)); err != nil {
		fmt.Println(err.Error())
	}
}
