package handler

import (
	"time"

	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/usecase"
	"auth-service/pkg/utils"
	"hospital-shared/jwt"
	"hospital-shared/logging"
	"hospital-shared/metrics"

	"auth-service/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	config      *config.Config
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		config:      cfg,
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
		metrics.RegisterFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	// Validation yap
	if err := middleware.ValidateStruct(&req); err != nil {
		metrics.RegisterFailCounter.Inc()
		var errors []string
		for _, validationErr := range err.(validator.ValidationErrors) {
			field := validationErr.Field()
			tag := validationErr.Tag()
			param := validationErr.Param()

			var message string
			switch tag {
			case "required":
				message = field + " is required"
			case "email":
				message = field + " must be a valid email"
			case "min":
				message = field + " must be at least " + param + " characters"
			case "max":
				message = field + " must be at most " + param + " characters"
			case "len":
				message = field + " must be exactly " + param + " characters"
			case "gt":
				message = field + " must be greater than " + param
			case "tc":
				message = field + " must be a valid TC identity number"
			case "phone":
				message = field + " must be a valid phone number"
			case "password":
				message = field + " must contain uppercase, lowercase and digit"
			default:
				message = field + " is invalid"
			}
			errors = append(errors, message)
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": errors,
		})
	}

	// Log registration attempt
	logging.GlobalLogger.LogInfo(
		c.UserContext(),
		"Registration attempt",
		zap.String("hospital_name", req.HospitalName),
		zap.String("authority_email", req.AuthorityEmail),
		zap.String("hospital_email", req.HospitalEmail),
	)

	authority, err := h.authUsecase.Register(&req)
	if err != nil {
		metrics.RegisterFailCounter.Inc()

		// Log registration failure
		logging.GlobalLogger.LogError(
			c.UserContext(),
			err,
			"Registration failed",
			zap.String("hospital_name", req.HospitalName),
			zap.String("authority_email", req.AuthorityEmail),
		)

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Log successful registration
	logging.GlobalLogger.LogInfo(
		c.UserContext(),
		"Registration successful",
		zap.String("hospital_name", req.HospitalName),
		zap.String("authority_email", req.AuthorityEmail),
		zap.Uint("authority_id", authority.ID),
		zap.Uint("hospital_id", authority.HospitalID),
	)

	metrics.RegisterSuccessCounter.Inc()

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
	metrics.LoginSuccessCounter.Inc()

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
		metrics.ForgotPasswordFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	resp, err := h.authUsecase.ForgotPassword(&req)
	if err != nil {
		metrics.ForgotPasswordFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	metrics.ForgotPasswordSuccessCounter.Inc()

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
		metrics.ResetPasswordFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.authUsecase.ResetPassword(&req); err != nil {
		metrics.ResetPasswordFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	metrics.ResetPasswordSuccessCounter.Inc()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password reset successful"})
}

// RefreshToken godoc
// @Summary     JWT yenileme
// @Description Geçerli bir refresh token ile yeni access ve refresh token döner
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       refreshToken body dto.RefreshTokenRequest true "Refresh token bilgisi"
// @Success     200 {object} dto.RefreshTokenResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /api/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		metrics.RefreshTokenFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if req.RefreshToken == "" {
		metrics.RefreshTokenFailCounter.Inc()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token is required",
		})
	}

	jwtCfg := utils.MapToSharedJWTConfig(h.config)
	tokenPair, err := jwt.RefreshAccessToken(req.RefreshToken, jwtCfg)
	if err != nil {
		metrics.RefreshTokenFailCounter.Inc()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}
	metrics.RefreshTokenSuccessCounter.Inc()

	return c.JSON(dto.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	})
}
