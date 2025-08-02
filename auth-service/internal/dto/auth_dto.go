package dto

import "time"

type RegisterRequest struct {
	HospitalName   string `json:"hospital_name"`
	TaxNumber      string `json:"tax_number"`
	HospitalEmail  string `json:"hospital_email"`
	HospitalPhone  string `json:"hospital_phone"`
	Address        string `json:"address"`
	CityID         uint   `json:"city_id"`
	DistrictID     uint   `json:"district_id"`
	AuthorityFName string `json:"authority_fname"`
	AuthorityLName string `json:"authority_lname"`
	AuthorityTC    string `json:"authority_tc"`
	AuthorityEmail string `json:"authority_email"`
	AuthorityPhone string `json:"authority_phone"`
	Password       string `json:"password"`
}

type LoginRequest struct {
	Credential string
	Password   string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ForgotPasswordRequest struct {
	Phone string `json:"phone"`
}

type ForgotPasswordResponse struct {
	Code string `json:"code"`
}

type ResetPasswordRequest struct {
	Phone             string `json:"phone"`
	Code              string `json:"code"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}

type AuthorityResponse struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	TC         string     `json:"tc"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Role       string     `json:"role"`
	HospitalID uint       `json:"hospital_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
