package main

import (
	"fmt"
	"log"

	"hospital-service/internal/config"
	"hospital-service/internal/database"
	"hospital-service/internal/router"
	"hospital-service/pkg/utils"

	"github.com/gofiber/fiber/v2"

	_ "hospital-service/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"hospital-shared/logging"
	"hospital-shared/metrics"
	"hospital-shared/tracing"
)

func main() {

	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Initialize structured logging
	logger := logging.InitLogger("hospital-service")
	logger.Info("Starting hospital service...")

	// Initialize distributed tracing
	tracerCloser, err := tracing.InitTracing("hospital-service")
	if err != nil {
		logger.Warn("Failed to initialize tracing: " + err.Error())
	} else {
		defer tracerCloser.Close()
		logger.Info("Distributed tracing initialized")
	}

	// Connect to database
	dbInstance, err := database.NewDatabase(&cfg)
	if err != nil {
		logger.Fatal("cannot connect to database: " + err.Error())
	}
	logger.Info("Database connection established")

	// Migration
	if err := database.RunMigrations(dbInstance.SQL); err != nil {
		logger.Fatal("migration failed: " + err.Error())
	}
	logger.Info("Database migrations completed")

	app := fiber.New()

	// Add observability middlewares
	app.Use(logging.CorrelationIDMiddleware())
	app.Use(logging.RequestLoggingMiddleware())
	app.Use(logging.ErrorLoggingMiddleware())

	app.Use(metrics.PrometheusMiddleware())

	app.Get("/metrics", metrics.PrometheusHandler())
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	deps := router.RouterDeps{
		App:             app,
		DB:              dbInstance,
		Config:          &cfg,
		JWTSharedConfig: utils.MapToSharedJWTConfig(&cfg),
	}

	router.HospitalRoutes(deps)
	router.PolyclinicRoutes(deps)
	router.LocationRoutes(deps)

	for _, r := range app.GetRoutes() {
		fmt.Println(r.Method, r.Path)
	}

	logger.Info(fmt.Sprintf("Hospital service starting on port %s", cfg.Server.Port))
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
