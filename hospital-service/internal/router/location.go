package router

import (
	"github.com/gofiber/fiber/v2"
	"hospital-service/internal/usecase"
	"hospital-service/internal/handler"
	"hospital-service/internal/repository"
)

func LocationRoutes(app *fiber.App) {

	locationRepo := repository.NewLocationRepository()
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	locationHandler := handler.NewLocationHandler(locationUsecase)

	api := app.Group("/api")

	locationGroup := api.Group("/location")

	locationGroup.Get("/cities", locationHandler.ListCities)
	locationGroup.Get("/districts", locationHandler.ListDistrictsByCity)
}
