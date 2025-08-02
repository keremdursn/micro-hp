package router

import (
	"hospital-service/internal/handler"
	"hospital-service/internal/infrastructure/client"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	"hospital-service/pkg/middleware"
	"hospital-service/pkg/utils"
)

func PolyclinicRoutes(deps RouterDeps) {

	polyclinicRepo := repository.NewPolyclinicRepository(deps.DB.SQL)
	personnelClient := client.NewPersonnelClient(deps.Config.Url.BaseUrl)
	polyclinicUsecase := usecase.NewPolyclinicUsecase(polyclinicRepo, personnelClient)
	polyclinicHandler := handler.NewPolyclinicHandler(polyclinicUsecase, deps.Config)

	api := deps.App.Group("/api")

	polyclinicGroup := api.Group("/polyclinic")

	polyclinicGroup.Get("/", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListAllPolyclinics)
	polyclinicGroup.Post("/hospital-polyclinics", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), polyclinicHandler.AddHospitalPolyclinic)
	polyclinicGroup.Get("/hospital-polyclinics", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.ListHospitalPolyclinic)
	polyclinicGroup.Delete("/hospital-polyclinics/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), polyclinicHandler.RemoveHospitalPolyclinic)

	//Personnel servisi bu endpointe http isteği atıyor
	polyclinicGroup.Get("/hospital-polyclinics/:id", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), polyclinicHandler.GetHospitalPolyclinic)
}
