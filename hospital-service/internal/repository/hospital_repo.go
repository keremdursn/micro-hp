package repository

import (
	"errors"

	"hospital-service/internal/models"

	"gorm.io/gorm"
)

type HospitalRepository interface {
	IsHospitalExists(taxNumber, email, phone string) (bool, error)
	CreateHospital(h *models.Hospital) error
	GetByID(id uint) (*models.Hospital, error)
	IsUniqueFieldsConflict(id uint, taxNumber, email, phone string) (bool, error)
	GetCityByID(id uint) (*models.City, error)
	GetDistrictByID(id uint) (*models.District, error)
	Update(hospital *models.Hospital) error
}

type hospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) HospitalRepository {
	return &hospitalRepository{db: db}
}

func (r *hospitalRepository) IsHospitalExists(taxNumber, email, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Hospital{}).
		Where("tax_number = ? OR email = ? OR phone = ?", taxNumber, email, phone).
		Count(&count).Error
	return count > 0, err
}

func (r *hospitalRepository) CreateHospital(h *models.Hospital) error {
	return r.db.Create(h).Error
}

func (r *hospitalRepository) GetByID(id uint) (*models.Hospital, error) {
	var hospital models.Hospital
	if err := r.db.First(&hospital, id).Error; err != nil {
		return nil, errors.New("hospital not found")
	}
	return &hospital, nil
}

func (r *hospitalRepository) IsUniqueFieldsConflict(id uint, taxNumber, email, phone string) (bool, error) {
	var count int64
	r.db.Model(&models.Hospital{}).
		Where("id != ? AND (tax_number = ? OR email = ? OR phone = ?)", id, taxNumber, email, phone).
		Count(&count)
	return count > 0, nil
}

func (r *hospitalRepository) GetCityByID(id uint) (*models.City, error) {
	var city models.City
	if err := r.db.First(&city, id).Error; err != nil {
		return nil, errors.New("city not found")
	}
	return &city, nil
}

func (r *hospitalRepository) GetDistrictByID(id uint) (*models.District, error) {
	var district models.District
	if err := r.db.First(&district, id).Error; err != nil {
		return nil, errors.New("district not found")
	}
	return &district, nil
}

func (r *hospitalRepository) Update(h *models.Hospital) error {
	return r.db.Save(h).Error
}
