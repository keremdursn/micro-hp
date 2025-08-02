package database

import (
	"fmt"

	"personnel-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	return db.AutoMigrate(
		&models.JobGroup{},
		&models.Title{},
		&models.Staff{},
	)
}
