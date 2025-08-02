package router

import (
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"auth-service/pkg/middleware"
	"auth-service/pkg/utils"
)

func SubUserRoutes(deps RouterDeps) {

	subuserRepo := repository.NewSubUserRepository(deps.DB.SQL)
	subuserUsecase := usecase.NewSubUserUsecase(subuserRepo)
	subuserHandler := handler.NewSubUserHandler(subuserUsecase, deps.Config)

	api := deps.App.Group("/api")

	subuserGroup := api.Group("/subuser")

	subuserGroup.Post("/", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), subuserHandler.CreateSubUser)
	subuserGroup.Get("/users", middleware.GeneralRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), subuserHandler.ListUsers)
	subuserGroup.Put("/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), subuserHandler.UpdateSubUser)
	subuserGroup.Delete("/:id", middleware.AdminRateLimiter(), utils.AuthRequired(deps.Config), utils.RequireRole("yetkili"), subuserHandler.DeleteSubUser)
}
