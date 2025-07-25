package repository

import (
	"hospital-service/internal/database"
	"hospital-service/internal/models"
)

type LocationRepository interface {
	GetAllCities() ([]models.City, error)
	GetDistrictsByCity(cityID uint) ([]models.District, error)
}

type locationRepository struct {
}

func NewLocationRepository() LocationRepository {
	return &locationRepository{}
}

func (r *locationRepository) GetAllCities() ([]models.City, error) {
	var cities []models.City
	if err := database.DB.Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *locationRepository) GetDistrictsByCity(cityID uint) ([]models.District, error) {
	var districts []models.District
	if err := database.DB.Where("city_id = ?", cityID).Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}
