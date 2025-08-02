package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/pkg/middleware"
)

func AuthRoutes(deps RouterDeps) {
	authRepo := repository.NewAuthRepository(deps.DB.SQL)
	authUsecase := usecase.NewAuthUsecase(authRepo, deps.DB.Redis)
	authHandler := handler.NewAuthHandler(authUsecase, deps.Config)

	api := deps.App.Group("/api")

	authGroup := api.Group("/auth")

	authGroup.Post("/register", middleware.AuthRateLimiter(), authHandler.Register)
	authGroup.Post("/login", middleware.LoginRateLimiter(), authHandler.Login)
	authGroup.Post("/forgot-password", middleware.AuthRateLimiter(), authHandler.ForgotPassword)
	authGroup.Post("/reset-password", middleware.AuthRateLimiter(), authHandler.ResetPassword)
	authGroup.Post("/refresh-token", middleware.AuthRateLimiter(), authHandler.RefreshToken)
}
