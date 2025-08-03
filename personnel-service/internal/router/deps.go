package router

import (
	sharedjwt "hospital-shared/jwt"
	"personnel-service/internal/config"
	"personnel-service/internal/database"

	"github.com/gofiber/fiber/v2"
)

type RouterDeps struct {
	App             *fiber.App
	DB              *database.Database
	Config          *config.Config
	JWTSharedConfig *sharedjwt.JWTConfig
}
