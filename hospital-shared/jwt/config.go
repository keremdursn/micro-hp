package jwt

type JWTConfig struct {
	PrivateKeyPEM      string
	PublicKeyPEM       string
	AccessTokenExpiry  string
	RefreshTokenExpiry string
}
