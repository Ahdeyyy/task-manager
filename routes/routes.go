package routes

import (
	"database/sql"
	"task-manager/store"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

type Routes struct {
	db      *sql.DB
	storage *postgres.Storage
	store   *session.Store
}

func NewApp(s *store.Store, connectionUri string) *fiber.App {
	// TODO: Pass the database into the app so it can be close properly
	app := fiber.New()

	storage := postgres.New(postgres.Config{
		Table:         "session_store",
		Reset:         false,
		GCInterval:    3 * time.Second,
		ConnectionURI: connectionUri,
	})
	session_store := session.New(session.Config{
		Storage:    storage,
		Expiration: 1 * time.Hour,
		KeyLookup:  "cookie:myapp_session",
	})

	routes := Routes{
		db:      s.Db,
		storage: storage,
		store:   session_store,
	}

	app.Static("/static/", "./static")
	app.Get("/", routes.HomeRoute)
	app.Get("/signup", routes.SignupRoute)
	app.Post("/users/create", routes.CreateNewUser)
	app.Post("/login", routes.LoginRoute)
	app.Get("/login", routes.LoginRoute)

	return app
}
