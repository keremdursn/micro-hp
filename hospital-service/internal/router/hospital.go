package router

import (
	"hospital-service/internal/database"
	"hospital-service/internal/handler"
	"hospital-service/internal/repository"
	"hospital-service/internal/usecase"
	utilss "hospital-shared/util"

	"github.com/gofiber/fiber/v2"
)

func HospitalRoutes(app *fiber.App, secret string) {
	db := database.GetDB()
	hRepo := repository.NewHospitalRepository(db)
	hUsecase := usecase.NewHospitalUsecase(hRepo)
	hHandler := handler.NewHospitalHandler(hUsecase)

	api := app.Group("/api")

	hGroup := api.Group("/hospital")

	hGroup.Get("/me", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), hHandler.GetHospitalMe)
	hGroup.Post("/", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), hHandler.GetHospitalMe)
	hGroup.Put("/me", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), hHandler.UpdateHospitalMe)
}
