package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/usecase"
	"auth-service/pkg/metrics"
	
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	config      *config.Config
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Register godoc
// @Summary     Hastane ve ilk yetkili kaydı
// @Description Registers a hospital and its first authority
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       register body dto.RegisterRequest true "Register info"
// @Success     201 {object} dto.AuthorityResponse
// @Failure     400 {object} map[string]string
// @Failure     409 {object} map[string]string
// @Router      /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	authority, err := h.authUsecase.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var deletedAt *time.Time
	if authority.DeletedAt.Valid {
		deletedAt = &authority.DeletedAt.Time
	}
	resp := dto.AuthorityResponse{
		ID:         authority.ID,
		FirstName:  authority.FirstName,
		LastName:   authority.LastName,
		TC:         authority.TC,
		Email:      authority.Email,
		Phone:      authority.Phone,
		Role:       authority.Role,
		HospitalID: authority.HospitalID,
		CreatedAt:  authority.CreatedAt,
		UpdatedAt:  authority.UpdatedAt,
		DeletedAt:  deletedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login godoc
// @Summary     Kullanıcı girişi
// @Description User login with email or phone, returns JWT token
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       login body dto.LoginRequest true "Login info"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	resp, err := h.authUsecase.Login(&req, h.config)
	if err != nil {
		// Başarısız login denemesi için Prometheus metriğini artır
		metrics.LoginFailCounter.Inc()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ForgotPassword godoc
// @Summary     Şifre sıfırlama kodu gönderir
// @Description Sends a reset code to the user's phone
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       forgot body dto.ForgotPasswordRequest true "Forgot password info"
// @Success     200 {object} dto.ForgotPasswordResponse
// @Failure     400 {object} map[string]string
// @Router      /api/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	resp, err := h.authUsecase.ForgotPassword(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ResetPassword godoc
// @Summary     Şifreyi sıfırlar
// @Description Resets the user's password with the code
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       reset body dto.ResetPasswordRequest true "Reset password info"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.authUsecase.ResetPassword(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successful"})
}
