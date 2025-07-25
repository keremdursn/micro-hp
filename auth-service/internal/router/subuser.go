package router

import (
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	utilss "hospital-shared/util"

	"github.com/gofiber/fiber/v2"
)

func SubUserRoutes(app *fiber.App, secret string) {
	db := database.GetDB()
	subuserRepo := repository.NewSubUserRepository(db)
	subuserUsecase := usecase.NewSubUserUsecase(subuserRepo)
	subuserHandler := handler.NewSubUserHandler(subuserUsecase)

	api := app.Group("/api")

	subuserGroup := api.Group("/subuser")

	subuserGroup.Post("/", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), subuserHandler.CreateSubUser)
	subuserGroup.Get("/subusers", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), subuserHandler.ListSubUsers)
	subuserGroup.Put("/:id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), subuserHandler.UpdateSubUser)
	subuserGroup.Delete("/:id", utilss.AuthRequired(secret), utilss.RequireRole("yetkili"), subuserHandler.DeleteSubUser)
}
