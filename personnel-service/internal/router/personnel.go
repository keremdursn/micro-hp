package router

import (
	utilss "hospital-shared/util"
	"personnel-service/internal/database"
	"personnel-service/internal/handler"
	"personnel-service/internal/repository"
	"personnel-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func PersonnelRoutes(app *fiber.App, secret string) {
	db := database.GetDB()
	personnelRepo := repository.NewPersonnelRepository(db)
	personnelUsecase := usecase.NewPersonnelUsecase(personnelRepo)
	personnelHandler := handler.NewPersonnelHandler(personnelUsecase)

	api := app.Group("/api")

	personnelGroup := api.Group("/personnel")

	personnelGroup.Get("/job-groups", utilss.AuthRequired(secret), utilss.RequireRole("yetkili", "calisan"), personnelHandler.ListAllJobGroups)
	personnelGroup.Get("/titles/:job_group_id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili", "calisan"), personnelHandler.ListTitleByJobGroup)

	personnelGroup.Post("/staff", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), personnelHandler.AddStaff)
	personnelGroup.Put("/staff/:id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), personnelHandler.UpdateStaff)
	personnelGroup.Delete("/staff/:id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), personnelHandler.DeleteStaff)
	personnelGroup.Get("/staff", utilss.AuthRequired(secret), utilss.RequireRole("yetkili", "calisan"), personnelHandler.ListStaff)
}
