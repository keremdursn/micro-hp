package router

import (
	"hospital-shared/jwt"
	"personnel-service/internal/handler"
	"personnel-service/internal/infrastructure/client"
	"personnel-service/internal/repository"
	"personnel-service/internal/usecase"

	"hospital-shared/middleware"
)

func PersonnelRoutes(deps RouterDeps) {
	personnelRepo := repository.NewPersonnelRepository(deps.DB.SQL)
	polyclinicClient := client.NewPolyclinicClient(deps.Config.Url.BaseUrl)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo, deps.DB.Redis, polyclinicClient)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, deps.Config)

	api := deps.App.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)

	//Hospital servisi bu endpointlere http istekleri atıyor
	personnelGroup.Get("/:id", jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), personnelHandler.GetStaffCount)
	personnelGroup.Get("/groups/:id", jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili", "calisan"), personnelHandler.GetGroupCounts)

}
