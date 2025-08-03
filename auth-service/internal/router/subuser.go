package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"hospital-shared/middleware"
	"hospital-shared/jwt"
)

func SubUserRoutes(deps RouterDeps) {

	subuserRepo := repository.NewSubUserRepository(deps.DB.SQL)
	subuserUsecase := usecase.NewSubUserUsecase(subuserRepo)
	subuserHandler := handler.NewSubUserHandler(subuserUsecase, deps.Config)

	api := deps.App.Group("/api")

	subuserGroup := api.Group("/subuser")

	subuserGroup.Post("/", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), subuserHandler.CreateSubUser)
	subuserGroup.Get("/users", middleware.GeneralRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), subuserHandler.ListUsers)
	subuserGroup.Put("/:id", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), subuserHandler.UpdateSubUser)
	subuserGroup.Delete("/:id", middleware.AdminRateLimiter(), jwt.AuthRequired(deps.JWTSharedConfig), jwt.RequireRole("yetkili"), subuserHandler.DeleteSubUser)
}
