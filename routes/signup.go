package routes

import (
	"errors"
	"task-manager/views"

	"github.com/gofiber/fiber/v2"
)

func (r Routes) SignupRoute(c *fiber.Ctx) error {
	return views.Render(c, views.Signup(""))
}

func (r Routes) CreateNewUser(c *fiber.Ctx) error {

	username := c.FormValue("username")
	password := c.FormValue("password") //TODO: use bcrypt to hash the password!!!
	email := c.FormValue("email")

	if username == "" || email == "" || password == "" {
		return errors.New("username or email or password cannot be empty")
	}

	statement := `
			INSERT INTO users (username, email, password) VALUES ($1,$2,$3)
		`

	_, err := r.db.Exec(statement, username, email, password)

	if err != nil {
		return views.Render(c, views.Signup("Username or email already exists"))
	}
	return c.Redirect("/login")

}
