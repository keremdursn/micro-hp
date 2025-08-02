package router

import (
	"github.com/gofiber/fiber/v2"
	"personnel-service/internal/config"
	"personnel-service/internal/database"
)

type RouterDeps struct {
	App    *fiber.App
	DB     *database.Database
	Config *config.Config
}
