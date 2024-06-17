package main

import (
	"fmt"
	"log"
	"task-manager/routes"
	"task-manager/store"
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	env_db := viper.GetString("DB_CONNECTION_URI")
	port := viper.GetString("PORT")

	if port == "" {
		port = "3000"
	}

	connectionUri := "postgres://postgres:password@localhost:5432/?sslmode=disable"
	if env_db != "" {
		connectionUri = env_db
		log.Println(connectionUri)
	}

	s, err := store.CreateStore(connectionUri)

	if err != nil {
		panic(err)
	}
	defer s.Db.Close()
	app := routes.NewApp(s, connectionUri)

	app.Use(logger.New())
	app.Get("/metrics", monitor.New())
	app.Use(limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	app.Listen(fmt.Sprintf(":%s", port))
}
