package router

import (
	"hospital-service/internal/handler"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	"hospital-shared/jwt"
	"hospital-shared/middleware"
)

func HospitalRoutes(deps RouterDeps) {
	hRepo := repository.NewHospitalRepository(deps.DB.SQL)
	hUsecase := usecase.NewHospitalUsecase(hRepo)
	hHandler := handler.NewHospitalHandler(hUsecase, deps.Config)

	api := deps.App.Group("/api")

	hGroup := api.Group("/hospital")

	hGroup.Get("/me", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), hHandler.GetHospitalMe)
	hGroup.Put("/me", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), hHandler.UpdateHospitalMe)
}
