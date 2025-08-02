package router

import (
	"github.com/gofiber/fiber/v2"
	"auth-service/internal/config"
	"auth-service/internal/database"
)

type RouterDeps struct {
	App    *fiber.App
	DB     *database.Database
	Config *config.Config
}
