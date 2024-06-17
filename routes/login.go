package routes

import (
	"errors"
	"task-manager/models"
	"task-manager/sessions"
	"task-manager/views"

	"github.com/gofiber/fiber/v2"
)

func (r Routes) LoginRoute(c *fiber.Ctx) error {

	email := c.FormValue("email")
	password := c.FormValue("password")

	if (email == "" || password == "") && c.Method() == "POST" {
		return errors.New("username or password cannot be empty")
	}

	statement := `
			SELECT * FROM users WHERE email=$1 AND password=$2
		`

	row := r.db.QueryRow(statement, email, password)

	user := models.User{}
	err := row.Scan(&user.Username, &user.Email, &user.Password)

	if err != nil && c.Method() == "POST" {
		return views.Render(c, views.Login("Invalid username or password", user.Email))
	}

	if c.Method() == "POST" {
		err = sessions.CreateUserSession(r.store, c, r.db, user.Username)
		if err != nil {
			return err
		}
	}
	if c.Method() == "POST" {
		return c.Redirect("/")
	}

	return views.Render(c, views.Login("", ""))

}
