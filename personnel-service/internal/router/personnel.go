package router

import (
	"personnel-service/internal/handler"
	"personnel-service/internal/repository"
	"personnel-service/internal/usecase"
	"personnel-service/pkg/utils"
	"personnel-service/internal/infrastructure/client"

	"personnel-service/pkg/middleware"
)

func PersonnelRoutes(deps RouterDeps) {
	personnelRepo := repository.NewPersonnelRepository(deps.DB.SQL)
	polyclinicClient := client.NewPolyclinicClient(deps.Config.Url.BaseUrl)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo, deps.DB.Redis, polyclinicClient)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase, deps.Config)

	api := deps.App.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)

	//Hospital servisi bu endpointlere http istekleri atıyor
	personnelGroup.Get("/:id", utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.GetStaffCount)
	personnelGroup.Get("/groups/:id", utils.AuthRequired(deps.Config), utils.RequireRole("yetkili", "calisan"), personnelHandler.GetGroupCounts)

}
