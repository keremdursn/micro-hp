package main

import (
	"fmt"
	"log"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/router"

	"github.com/gofiber/fiber/v2"

	// _ "auth-service/docs"

	// fiberSwagger "github.com/swaggo/fiber-swagger"

	"auth-service/pkg/metrics"
)

func main() {

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Connect to database
	dbInstance, err := database.NewDatabase(&cfg)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// Migration
	if err := database.RunMigrations(dbInstance.SQL); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err != nil {
		log.Fatal("cannot migrate database: ", err)
	}

	app := fiber.New()

	app.Use(metrics.PrometheusMiddleware())
	app.Get("/metrics", metrics.PrometheusHandler())

	// app.Get("/swagger/*", fiberSwagger.WrapHandler)

	deps := router.RouterDeps{
		App:    app,
		DB:     dbInstance,
		Config: &cfg,
	}

	router.AuthRoutes(deps)
	router.SubUserRoutes(deps)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
