package routes

import (
	"task-manager/views"

	"github.com/gofiber/fiber/v2"
)

func HomeRoute(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		name = "World"
	}
	return views.Render(c, views.Home(name))
}

func Clicked(c *fiber.Ctx) error {
	return views.Render(c, views.Clicked())
}
