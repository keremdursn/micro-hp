package main

import (
	"fmt"
	"log"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/models"
	"auth-service/internal/router"

	"github.com/gofiber/fiber/v2"

	_ "auth-service/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"auth-service/pkg/metrics"
)

func main() {

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	secret := cfg.JWT.Secret

	// Connect to database
	database.Connect(&cfg)
	database.ConnectRedis(&cfg)

	err = database.DB.AutoMigrate(

		&models.Authority{},
	)
	if err != nil {
		log.Fatal("cannot migrate database: ", err)
	}

	app := fiber.New()

	app.Use(metrics.PrometheusMiddleware())
	app.Get("/metrics", metrics.PrometheusHandler())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	router.AuthRoutes(app)
	router.SubUserRoutes(app, secret)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
