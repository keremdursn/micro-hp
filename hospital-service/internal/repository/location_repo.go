package repository

import (
	"hospital-service/internal/models"

	"gorm.io/gorm"
)

type LocationRepository interface {
	GetAllCities() ([]models.City, error)
	GetDistrictsByCity(cityID uint) ([]models.District, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{
		db: db,
	}
}

func (r *locationRepository) GetAllCities() ([]models.City, error) {
	var cities []models.City
	if err := r.db.Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *locationRepository) GetDistrictsByCity(cityID uint) ([]models.District, error) {
	var districts []models.District
	if err := r.db.Where("city_id = ?", cityID).Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}
