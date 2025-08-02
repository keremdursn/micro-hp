package middleware

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct DTO validation için basit fonksiyon
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Custom validation functions
func init() {
	// TC Kimlik validation
	validate.RegisterValidation("tc", validateTC)

	// Phone validation
	validate.RegisterValidation("phone", validatePhone)

	// Password strength validation
	validate.RegisterValidation("password", validatePassword)
}

// validateTC TC kimlik numarası kontrolü
func validateTC(fl validator.FieldLevel) bool {
	tc := fl.Field().String()

	// 11 haneli olmalı
	if len(tc) != 11 {
		return false
	}

	// Sadece rakam olmalı
	for _, char := range tc {
		if char < '0' || char > '9' {
			return false
		}
	}

	// İlk hane 0 olamaz
	if tc[0] == '0' {
		return false
	}

	// TC kimlik algoritması kontrolü
	// Gerçek uygulamada daha detaylı algoritma kullanılabilir
	return true
}

// validatePhone telefon numarası kontrolü
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Türkiye telefon formatı kontrolü
	// 5xx xxx xx xx veya +90 5xx xxx xx xx
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// +90 ile başlıyorsa kaldır
	if strings.HasPrefix(phone, "+90") {
		phone = phone[3:]
	}

	// 0 ile başlıyorsa kaldır
	if strings.HasPrefix(phone, "0") {
		phone = phone[1:]
	}

	// 10 haneli olmalı ve 5 ile başlamalı
	return len(phone) == 10 && phone[0] == '5'
}

// validatePassword şifre gücü kontrolü
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// En az 6 karakter
	if len(password) < 6 {
		return false
	}

	// En az bir büyük harf
	hasUpper := false
	// En az bir küçük harf
	hasLower := false
	// En az bir rakam
	hasDigit := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
