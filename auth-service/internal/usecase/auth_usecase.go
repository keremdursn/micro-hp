package usecase

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/dto"
	"auth-service/internal/infrastructure/client"
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/utils"
	dt "hospital-shared/dto"
	"hospital-shared/jwt"
	"hospital-shared/logging"
	"hospital-shared/tracing"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *dto.RegisterRequest) (*models.Authority, error)
	Login(req *dto.LoginRequest, cfg *config.Config) (*dto.LoginResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

type authUsecase struct {
	authRepo       repository.AuthRepository
	redis          *redis.Client
	hospitalClient client.HospitalClient
}

func NewAuthUsecase(r repository.AuthRepository, redis *redis.Client, hc client.HospitalClient) AuthUsecase {
	return &authUsecase{
		authRepo:       r,
		redis:          redis,
		hospitalClient: hc,
	}
}

func (u *authUsecase) Register(req *dto.RegisterRequest) (*models.Authority, error) {
	ctx := context.Background()

	// Start tracing span for registration process
	span, ctx := tracing.StartServiceSpan(ctx, "auth-service", "register")
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		tracing.FinishSpanWithError(span, err)
		logging.GlobalLogger.LogError(ctx, err, "Password hashing failed")
		return nil, err
	}

	// Check uniqueness with database tracing
	dbSpan, dbCtx := tracing.StartDatabaseSpan(ctx, "SELECT", "authorities")
	exists, err := u.authRepo.IsAuthorityExists(req.AuthorityTC, req.AuthorityEmail, req.AuthorityPhone)
	if err != nil {
		tracing.FinishSpanWithError(dbSpan, err)
		logging.GlobalLogger.LogDatabaseOperation(dbCtx, "SELECT", "authorities", 0, err)
		return nil, err
	}
	if dbSpan != nil {
		dbSpan.Finish()
	}
	logging.GlobalLogger.LogDatabaseOperation(dbCtx, "SELECT", "authorities", 0, nil)

	if exists {
		err := errors.New("authority already exists")
		tracing.FinishSpanWithError(span, err)
		logging.GlobalLogger.LogInfo(ctx, "Registration failed - authority exists",
			zap.String("tc", req.AuthorityTC),
			zap.String("email", req.AuthorityEmail),
		)
		return nil, err
	}

	// Call hospital service to create hospital with service tracing
	serviceSpan, serviceCtx := tracing.StartHTTPSpan(ctx, "create-hospital", "POST", "/api/hospital")
	start := time.Now()
	hospitalResp, err := u.hospitalClient.CreateHospital(&dt.CreateHospitalRequest{
		Name:       req.HospitalName,
		TaxNumber:  req.TaxNumber,
		Email:      req.HospitalEmail,
		Phone:      req.HospitalPhone,
		Address:    req.Address,
		CityID:     req.CityID,
		DistrictID: req.DistrictID,
	})
	duration := time.Since(start)

	if err != nil {
		tracing.FinishSpanWithError(serviceSpan, err)
		logging.GlobalLogger.LogServiceCall(serviceCtx, "hospital-service", "/api/hospital", 0, duration, err)
		return nil, err
	}
	if serviceSpan != nil {
		serviceSpan.Finish()
	}
	logging.GlobalLogger.LogServiceCall(serviceCtx, "hospital-service", "/api/hospital", 201, duration, nil)

	// Create Authority with database tracing
	dbSpan2, dbCtx2 := tracing.StartDatabaseSpan(ctx, "INSERT", "authorities")
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

	start2 := time.Now()
	err = u.authRepo.CreateAuthority(authority)
	duration2 := time.Since(start2)

	if err != nil {
		tracing.FinishSpanWithError(dbSpan2, err)
		logging.GlobalLogger.LogDatabaseOperation(dbCtx2, "INSERT", "authorities", duration2, err)
		return nil, err
	}
	if dbSpan2 != nil {
		dbSpan2.Finish()
	}
	logging.GlobalLogger.LogDatabaseOperation(dbCtx2, "INSERT", "authorities", duration2, nil)

	logging.GlobalLogger.LogInfo(ctx, "Registration completed successfully",
		zap.Uint("authority_id", authority.ID),
		zap.Uint("hospital_id", authority.HospitalID),
	)

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
	jwtCfg := utils.MapToSharedJWTConfig(cfg)
	tokenPair, err := jwt.GenerateTokenPair(authority.ID, authority.HospitalID, authority.Role, jwtCfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
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
	if err := u.redis.Set(ctx, "reset_code:"+req.Phone, code, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponse{Code: code}, nil
}

func (u *authUsecase) ResetPassword(req *dto.ResetPasswordRequest) error {
	if req.NewPassword != req.RepeatNewPassword {
		return errors.New("passwords do not match")
	}

	ctx := context.Background()
	storedCode, err := u.redis.Get(ctx, "reset_code:"+req.Phone).Result()
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

	_ = u.redis.Del(ctx, "reset_code:"+req.Phone).Err()
	return nil
}
