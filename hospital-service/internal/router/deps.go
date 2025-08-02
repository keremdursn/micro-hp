package router

import (
	"github.com/gofiber/fiber/v2"
	"hospital-service/internal/config"
	"hospital-service/internal/database"
)

type RouterDeps struct {
	App    *fiber.App
	DB     *database.Database
	Config *config.Config
}
