package main

import (
	"fmt"
	"log"

	"hospital-service/internal/config"
	"hospital-service/internal/database"
	"hospital-service/internal/models"
	"hospital-service/internal/router"

	"github.com/gofiber/fiber/v2"

	_ "auth-service/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"hospital-service/pkg/metrics"
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

		&models.Hospital{},
		&models.City{},
		&models.District{},
		&models.Polyclinic{},
		&models.HospitalPolyclinic{},
	)
	if err != nil {
		log.Fatal("cannot migrate database: ", err)
	}

	app := fiber.New()

	app.Use(metrics.PrometheusMiddleware())
	app.Get("/metrics", metrics.PrometheusHandler())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	router.HospitalRoutes(app, secret)
	router.PolyclinicRoutes(app, secret)
	router.LocationRoutes(app)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
