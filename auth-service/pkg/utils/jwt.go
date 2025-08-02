package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"auth-service/internal/config"
)

type UserInfo struct {
	AuthorityID uint
	HospitalID  uint
	Role        string
}

// UserContextKey is the key for user info in Fiber context
const UserContextKey = "userInfo"

// TokenPair represents both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // access token expiry in seconds
}

// Claims for access token
type AccessTokenClaims struct {
	AuthorityID uint   `json:"authority_id"`
	HospitalID  uint   `json:"hospital_id"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// Claims for refresh token
type RefreshTokenClaims struct {
	AuthorityID uint   `json:"authority_id"`
	HospitalID  uint   `json:"hospital_id"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

// LoadPrivateKey loads RSA private key from config string
func LoadPrivateKeyFromConfig(cfg *config.Config) (*rsa.PrivateKey, error) {
	keyBytes := []byte(cfg.JWT.PrivateKey)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from config private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return privateKey, nil
}

// LoadPublicKey loads RSA public key from config string
func LoadPublicKeyFromConfig(cfg *config.Config) (*rsa.PublicKey, error) {
	keyBytes := []byte(cfg.JWT.PublicKey)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from config public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPublicKey, nil
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(authorityID, hospitalID uint, role string, cfg *config.Config) (*TokenPair, error) {
	privateKey, err := LoadPrivateKeyFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	// Parse expiry durations
	accessExpiry, err := time.ParseDuration(cfg.JWT.AccessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("invalid access token expiry: %w", err)
	}

	refreshExpiry, err := time.ParseDuration(cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token expiry: %w", err)
	}

	// Create access token claims
	accessClaims := &AccessTokenClaims{
		AuthorityID: authorityID,
		HospitalID:  hospitalID,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hospital-api",
		},
	}

	// Create refresh token claims
	refreshClaims := &RefreshTokenClaims{
		AuthorityID: authorityID,
		HospitalID:  hospitalID,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hospital-api",
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessExpiry.Seconds()),
	}, nil
}

// ParseAccessToken validates and parses an access token
func ParseAccessToken(tokenStr string, cfg *config.Config) (*UserInfo, error) {
	publicKey, err := LoadPublicKeyFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid access token claims")
	}

	return &UserInfo{
		AuthorityID: claims.AuthorityID,
		HospitalID:  claims.HospitalID,
		Role:        claims.Role,
	}, nil
}

// ParseRefreshToken validates and parses a refresh token
func ParseRefreshToken(tokenStr string, cfg *config.Config) (*UserInfo, error) {
	publicKey, err := LoadPublicKeyFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return &UserInfo{
		AuthorityID: claims.AuthorityID,
		HospitalID:  claims.HospitalID,
		Role:        claims.Role,
	}, nil
}

// RefreshAccessToken creates new token pair using refresh token
func RefreshAccessToken(refreshToken string, cfg *config.Config) (*TokenPair, error) {
	// Parse refresh token to get user info
	userInfo, err := ParseRefreshToken(refreshToken, cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return GenerateTokenPair(userInfo.AuthorityID, userInfo.HospitalID, userInfo.Role, cfg)
}

// AuthRequired returns a Fiber middleware that checks JWT and sets user info in context
func AuthRequired(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		user, err := ParseAccessToken(tokenStr, cfg)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired access token",
			})
		}

		c.Locals(UserContextKey, user)
		return c.Next()
	}
}

// GetUserInfo extracts user info from Fiber context
func GetUserInfo(c *fiber.Ctx) *UserInfo {
	val := c.Locals(UserContextKey)
	if user, ok := val.(*UserInfo); ok {
		return user
	}
	return nil
}

// RequireRole returns a Fiber middleware that allows only users with the given role(s)
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := GetUserInfo(c)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: insufficient role"})
	}
}
