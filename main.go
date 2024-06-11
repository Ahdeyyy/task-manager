package main

import (
	"task-manager/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/", "./static")
	app.Get("/", routes.HomeRoute)
	app.Get("/name/:name?", routes.HomeRoute)
	app.Get("/clicked", routes.Clicked)

	app.Listen(":3000")
}
