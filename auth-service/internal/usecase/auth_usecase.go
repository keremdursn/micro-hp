package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/utils"
	dt "hospital-shared/dto"
	utilss "hospital-shared/util"

	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *dto.RegisterRequest) (*models.Authority, error)
	Login(req *dto.LoginRequest, cfg *config.Config) (*dto.LoginResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

type authUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(r repository.AuthRepository) AuthUsecase {
	return &authUsecase{authRepo: r}
}

func (u *authUsecase) Register(req *dto.RegisterRequest) (*models.Authority, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Check uniqueness
	if exists, _ := u.authRepo.IsAuthorityExists(req.AuthorityTC, req.AuthorityEmail, req.AuthorityPhone); exists {
		return nil, errors.New("authority already exists")
	}

	// Create Hospital
	hospitalReq := dt.CreateHospitalRequest{
		Name:       req.HospitalName,
		TaxNumber:  req.TaxNumber,
		Email:      req.HospitalEmail,
		Phone:      req.HospitalPhone,
		Address:    req.Address,
		CityID:     req.CityID,
		DistrictID: req.DistrictID,
	}
	body, _ := json.Marshal(hospitalReq)
	resp, err := http.Post("http://hospital-service:8082/api/hospital", "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusCreated {
		// hata yönetimi
	}
	var hospitalResp dt.HospitalResponse
	json.NewDecoder(resp.Body).Decode(&hospitalResp)

	// Create Authority
	authority := &models.Authority{
		FirstName:  req.AuthorityFName,
		LastName:   req.AuthorityLName,
		TC:         req.AuthorityTC,
		Email:      req.AuthorityEmail,
		Phone:      req.AuthorityPhone,
		Password:   hashedPassword,
		Role:       "yetkili",
		HospitalID: hospitalResp.ID,
	}
	if err := u.authRepo.CreateAuthority(authority); err != nil {
		return nil, err
	}

	return authority, nil
}

func (u *authUsecase) Login(req *dto.LoginRequest, cfg *config.Config) (*dto.LoginResponse, error) {
	authority, err := u.authRepo.GetAuthorityByEmailOrPhone(req.Credential)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Parola kontrolü
	if !utils.CheckPasswordHash(req.Password, authority.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Token üret
	secret := cfg.JWT.Secret
	token, err := utilss.GenerateToken(authority.ID, authority.HospitalID, authority.Role, secret)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (u *authUsecase) ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	_, err := u.authRepo.GetAuthorityByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	code := utils.GenerateResetCode()

	ctx := context.Background()
	if err := database.RDB.Set(ctx, "reset_code:"+req.Phone, code, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponse{Code: code}, nil
}

func (u *authUsecase) ResetPassword(req *dto.ResetPasswordRequest) error {
	if req.NewPassword != req.RepeatNewPassword {
		return errors.New("passwords do not match")
	}

	ctx := context.Background()
	storedCode, err := database.RDB.Get(ctx, "reset_code:"+req.Phone).Result()
	if err != nil || storedCode != req.Code {
		return errors.New("invalid or expired code")
	}

	authority, err := u.authRepo.GetAuthorityByPhone(req.Phone)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := u.authRepo.UpdateAuthorityPassword(authority, hashedPassword); err != nil {
		return err
	}

	_ = database.RDB.Del(ctx, "reset_code:"+req.Phone).Err()
	return nil
}
