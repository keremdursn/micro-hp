package main

import (
	"fmt"
	"log"

	"personnel-service/internal/config"
	"personnel-service/internal/database"
	"personnel-service/internal/models"
	"personnel-service/internal/router"

	"github.com/gofiber/fiber/v2"

	_ "personnel-service/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"personnel-service/pkg/metrics"
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

		&models.JobGroup{},
		&models.Title{},
		&models.Staff{},
	)
	if err != nil {
		log.Fatal("cannot migrate database: ", err)
	}

	app := fiber.New()

	app.Use(metrics.PrometheusMiddleware())
	app.Get("/metrics", metrics.PrometheusHandler())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	router.PersonnelRoutes(app, secret)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
