package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/infrastructure/client"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"hospital-shared/middleware"
)

func AuthRoutes(deps RouterDeps) {
	authRepo := repository.NewAuthRepository(deps.DB.SQL)
	hospitalClient := client.NewHospitalClient(deps.Config.Url.BaseUrl)
	authUsecase := usecase.NewAuthUsecase(authRepo, deps.DB.Redis, hospitalClient)
	authHandler := handler.NewAuthHandler(authUsecase, deps.Config)

	api := deps.App.Group("/api")

	authGroup := api.Group("/auth")

	authGroup.Post("/register", middleware.AuthRateLimiter(), authHandler.Register)
	authGroup.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	authGroup.Post("/forgot-password", middleware.AuthRateLimiter(), authHandler.ForgotPassword)
	authGroup.Post("/reset-password", middleware.AuthRateLimiter(), authHandler.ResetPassword)
	authGroup.Post("/refresh-token", middleware.AuthRateLimiter(), authHandler.RefreshToken)
}
