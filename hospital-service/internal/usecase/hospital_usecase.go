package usecase

import (
	"context"
	"errors"
	"time"

	"hospital-service/internal/dto"
	"hospital-service/internal/models"
	"hospital-service/internal/repository"
	dt "hospital-shared/dto"
	"hospital-shared/logging"
	"hospital-shared/tracing"

	"go.uber.org/zap"
)

type HospitalUsecase interface {
	GetHospitalByID(hospitalID uint) (*dto.HospitalResponse, error)
	UpdateHospital(hospitalID uint, req *dto.UpdateHospitalRequest) (*dto.HospitalResponse, error)
	CreateHospital(req *dt.CreateHospitalRequest) (*dto.HospitalResponse, error)
}

type hospitalUsecase struct {
	repo repository.HospitalRepository
}

func NewHospitalUsecase(repo repository.HospitalRepository) HospitalUsecase {
	return &hospitalUsecase{repo: repo}
}

func (u *hospitalUsecase) CreateHospital(req *dt.CreateHospitalRequest) (*dto.HospitalResponse, error) {
	ctx := context.Background()

	// Start tracing span for hospital creation
	span, ctx := tracing.StartServiceSpan(ctx, "hospital-service", "create-hospital")
	defer func() {
		if span != nil {
			span.Finish()
		}
	}()

	// Log hospital creation attempt
	logging.GlobalLogger.LogInfo(ctx, "Hospital creation attempt",
		zap.String("hospital_name", req.Name),
		zap.String("tax_number", req.TaxNumber),
	)

	// Check if hospital exists with database tracing
	dbSpan, dbCtx := tracing.StartDatabaseSpan(ctx, "SELECT", "hospitals")
	start := time.Now()
	exists, _ := u.repo.IsHospitalExists(req.TaxNumber, req.Email, req.Phone)
	duration := time.Since(start)
	tracing.FinishSpanWithError(dbSpan, nil)
	logging.GlobalLogger.LogDatabaseOperation(dbCtx, "SELECT", "hospitals", duration, nil)

	if exists {
		err := errors.New("hospital already exists")
		tracing.FinishSpanWithError(span, err)
		logging.GlobalLogger.LogError(ctx, err, "Hospital creation failed - already exists")
		return nil, err
	}

	hospital := &models.Hospital{
		Name:       req.Name,
		TaxNumber:  req.TaxNumber,
		Email:      req.Email,
		Phone:      req.Phone,
		Address:    req.Address,
		CityID:     req.CityID,
		DistrictID: req.DistrictID,
	}

	// Create hospital with database tracing
	dbSpan2, dbCtx2 := tracing.StartDatabaseSpan(ctx, "INSERT", "hospitals")
	start2 := time.Now()
	err := u.repo.CreateHospital(hospital)
	duration2 := time.Since(start2)
	tracing.FinishSpanWithError(dbSpan2, err)
	logging.GlobalLogger.LogDatabaseOperation(dbCtx2, "INSERT", "hospitals", duration2, err)

	if err != nil {
		tracing.FinishSpanWithError(span, err)
		logging.GlobalLogger.LogError(ctx, err, "Hospital creation failed")
		return nil, err
	}

	city, _ := u.repo.GetCityByID(hospital.CityID)
	district, _ := u.repo.GetDistrictByID(hospital.DistrictID)

	// Log successful creation
	logging.GlobalLogger.LogInfo(ctx, "Hospital created successfully",
		zap.Uint("hospital_id", hospital.ID),
		zap.String("hospital_name", hospital.Name),
	)

	return &dto.HospitalResponse{
		ID:           hospital.ID,
		TaxNumber:    hospital.TaxNumber,
		Name:         hospital.Name,
		Email:        hospital.Email,
		Phone:        hospital.Phone,
		Address:      hospital.Address,
		CityID:       hospital.CityID,
		CityName:     city.Name,
		DistrictID:   hospital.DistrictID,
		DistrictName: district.Name,
	}, nil
}

func (u *hospitalUsecase) GetHospitalByID(hospitalID uint) (*dto.HospitalResponse, error) {
	hospital, err := u.repo.GetByID(hospitalID)
	if err != nil {
		return nil, err
	}
	city, _ := u.repo.GetCityByID(hospital.CityID)
	district, _ := u.repo.GetDistrictByID(hospital.DistrictID)

	return &dto.HospitalResponse{
		ID:           hospital.ID,
		Name:         hospital.Name,
		TaxNumber:    hospital.TaxNumber,
		Email:        hospital.Email,
		Phone:        hospital.Phone,
		Address:      hospital.Address,
		CityID:       hospital.CityID,
		CityName:     city.Name,
		DistrictID:   hospital.DistrictID,
		DistrictName: district.Name,
	}, nil
}

func (u *hospitalUsecase) UpdateHospital(hospitalID uint, req *dto.UpdateHospitalRequest) (*dto.HospitalResponse, error) {
	hospital, err := u.repo.GetByID(hospitalID)
	if err != nil {
		return nil, err
	}
	conflict, _ := u.repo.IsUniqueFieldsConflict(hospitalID, req.TaxNumber, req.Email, req.Phone)
	if conflict {
		return nil, errors.New("another hospital with given tax number, email, or phone already exists")
	}

	city, err := u.repo.GetCityByID(req.CityID)
	if err != nil {
		return nil, err
	}
	district, err := u.repo.GetDistrictByID(req.DistrictID)
	if err != nil {
		return nil, err
	}

	hospital.Name = req.Name
	hospital.TaxNumber = req.TaxNumber
	hospital.Email = req.Email
	hospital.Phone = req.Phone
	hospital.Address = req.Address
	hospital.CityID = req.CityID
	hospital.DistrictID = req.DistrictID

	if err := u.repo.Update(hospital); err != nil {
		return nil, err
	}

	return &dto.HospitalResponse{
		ID:           hospital.ID,
		Name:         hospital.Name,
		TaxNumber:    hospital.TaxNumber,
		Email:        hospital.Email,
		Phone:        hospital.Phone,
		Address:      hospital.Address,
		CityID:       hospital.CityID,
		CityName:     city.Name,
		DistrictID:   hospital.DistrictID,
		DistrictName: district.Name,
	}, nil
}
