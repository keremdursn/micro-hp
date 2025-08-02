package database

import (
	"fmt"

	"auth-service/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	return db.AutoMigrate(

		&models.Authority{},
		
	)
}
