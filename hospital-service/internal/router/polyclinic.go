package router

import (
	"hospital-service/internal/database"
	"hospital-service/internal/handler"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	utilss "hospital-shared/util"

	"github.com/gofiber/fiber/v2"
)

func PolyclinicRoutes(app *fiber.App, secret string) {
	db := database.GetDB()
	polyclinicRepo := repository.NewPolyclinicRepository(db)
	polyclinicUsecase := usecase.NewPolyclinicUsecase(polyclinicRepo)
	polyclinicHandler := handler.NewPolyclinicHandler(polyclinicUsecase)

	api := app.Group("/api")

	polyclinicGroup := api.Group("/polyclinic")

	//
	polyclinicGroup.Get("/", utilss.AuthRequired(secret), utilss.RequireRole("yetkili", "calisan"), polyclinicHandler.ListAllPolyclinics)
	polyclinicGroup.Post("/hospital-polyclinics", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), polyclinicHandler.AddHospitalPolyclinic)
	polyclinicGroup.Get("/hospital-polyclinics", utilss.AuthRequired(secret), utilss.RequireRole("yetkili", "calisan"), polyclinicHandler.ListHospitalPolyclinic)
	polyclinicGroup.Delete("/hospital-polyclinics/:id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), polyclinicHandler.RemoveHospitalPolyclinic)
}
