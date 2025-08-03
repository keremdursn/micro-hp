package utils

import (
	sharedjwt "hospital-shared/jwt"
	"personnel-service/internal/config"
)

// Dönüştürücü fonksiyon
func MapToSharedJWTConfig(cfg *config.Config) *sharedjwt.JWTConfig {
	return &sharedjwt.JWTConfig{
		PrivateKeyPEM:      cfg.JWT.PrivateKey,
		PublicKeyPEM:       cfg.JWT.PublicKey,
		AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
	}
}
