package router

import (
	"github.com/gofiber/fiber/v2"
	"hospital-service/internal/config"
	"hospital-service/internal/database"
	sharedjwt "hospital-shared/jwt"
)

type RouterDeps struct {
	App    *fiber.App
	DB     *database.Database
	Config *config.Config
	JWTSharedConfig *sharedjwt.JWTConfig
}
