package router

import (
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	db := database.GetDB()
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	api := app.Group("/api")

	authGroup := api.Group("/auth")

	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/forgot-password", authHandler.ForgotPassword)
	authGroup.Post("/reset-password", authHandler.ResetPassword)

	// Sub-user management (only for 'yetkili')
	// api.Post("/users", utils.AuthRequired(cfg), utils.RequireRole("yetkili"), authHandler.CreateSubUser)
}
