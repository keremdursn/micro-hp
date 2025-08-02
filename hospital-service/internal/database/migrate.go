package database

import (
	"fmt"

	"hospital-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	return db.AutoMigrate(
		&models.Hospital{},
		&models.City{},
		&models.District{},
		&models.Polyclinic{},
		&models.HospitalPolyclinic{},
	)
}
