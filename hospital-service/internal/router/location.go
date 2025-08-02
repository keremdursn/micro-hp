package router

import (
	"hospital-service/internal/handler"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	"hospital-service/pkg/middleware"
)

func LocationRoutes(deps RouterDeps) {

	locationRepo := repository.NewLocationRepository(deps.DB.SQL)
	locationUsecase := usecase.NewLocationUsecase(locationRepo, deps.DB.Redis)
	locationHandler := handler.NewLocationHandler(locationUsecase)

	api := deps.App.Group("/api")

	locationGroup := api.Group("/location")

	locationGroup.Get("/cities", middleware.GeneralRateLimiter(), locationHandler.ListCities)
	locationGroup.Get("/districts", middleware.GeneralRateLimiter(), locationHandler.ListDistrictsByCity)
}
