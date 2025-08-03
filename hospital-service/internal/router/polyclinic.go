package router

import (
	"hospital-service/internal/handler"
	"hospital-service/internal/infrastructure/client"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	"hospital-shared/jwt"
	"hospital-shared/middleware"
)

func PolyclinicRoutes(deps RouterDeps) {

	polyclinicRepo := repository.NewPolyclinicRepository(deps.DB.SQL)
	personnelClient := client.NewPersonnelClient(deps.Config.Url.BaseUrl)
	polyclinicUsecase := usecase.NewPolyclinicUsecase(polyclinicRepo, personnelClient)
	polyclinicHandler := handler.NewPolyclinicHandler(polyclinicUsecase, deps.Config)

	api := deps.App.Group("/api")

	polyclinicGroup := api.Group("/polyclinic")

	polyclinicGroup.Get("/", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), polyclinicHandler.ListAllPolyclinics)
	polyclinicGroup.Post("/hospital-polyclinics", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), polyclinicHandler.AddHospitalPolyclinic)
	polyclinicGroup.Get("/hospital-polyclinics", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), polyclinicHandler.ListHospitalPolyclinic)
	polyclinicGroup.Delete("/hospital-polyclinics/:id", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), polyclinicHandler.RemoveHospitalPolyclinic)

	//Personnel servisi bu endpointe http isteği atıyor
	polyclinicGroup.Get("/hospital-polyclinics/:id", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), polyclinicHandler.GetHospitalPolyclinic)
}
