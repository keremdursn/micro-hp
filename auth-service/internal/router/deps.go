package router

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	sharedjwt "hospital-shared/jwt"

	"github.com/gofiber/fiber/v2"
)

type RouterDeps struct {
	App             *fiber.App
	DB              *database.Database
	Config          *config.Config
	JWTSharedConfig *sharedjwt.JWTConfig
}
