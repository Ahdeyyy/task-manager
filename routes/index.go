package routes

import (
	"log"
	"task-manager/sessions"
	"task-manager/views"

	"github.com/gofiber/fiber/v2"
)

func (r Routes) HomeRoute(c *fiber.Ctx) error {

	user, err := sessions.GetUserSessionData(r.db, r.store, c)
	log.Println(user)
	if err != nil {
		log.Println(err)
		return views.Render(c, views.Home(nil))
	}
	return views.Render(c, views.Home(user))

}
