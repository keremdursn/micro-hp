package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type UserInfo struct {
	AuthorityID uint
	HospitalID  uint
	Role        string
}

// UserContextKey is the key for user info in Fiber context
const UserContextKey = "userInfo"

type Claims struct {
	AuthorityID uint   `json:"authority_id"`
	HospitalID  uint   `json:"hospital_id"`
	Role        string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(authorityID, hospitalID uint, role string, secret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		AuthorityID: authorityID,
		HospitalID:  hospitalID,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates and parses a JWT, returning user info
func ParseToken(tokenStr string, secret string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return &UserInfo{
		AuthorityID: claims.AuthorityID,
		HospitalID:  claims.HospitalID,
		Role:        claims.Role,
	}, nil
}

// AuthRequired returns a Fiber middleware that checks JWT and sets user info in context
func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		user, err := ParseToken(tokenStr, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
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
