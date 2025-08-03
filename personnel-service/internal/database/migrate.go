package database

import (
	"fmt"

	"personnel-service/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	err := db.AutoMigrate(
		&models.JobGroup{},
		&models.Title{},
		&models.Staff{},
	)
	if err != nil {
		return err
	}

	// Seed data ekleme
	return seedData(db)
}

func seedData(db *gorm.DB) error {
	// JobGroup seed data ekleme
	var jobGroupCount int64
	db.Model(&models.JobGroup{}).Count(&jobGroupCount)

	if jobGroupCount == 0 {
		jobGroups := []models.JobGroup{
			{Name: "Doktor"},
			{Name: "Hemşire"},
			{Name: "Teknisyen"},
			{Name: "İdari Personel"},
			{Name: "Güvenlik"},
		}

		for _, jobGroup := range jobGroups {
			if err := db.Create(&jobGroup).Error; err != nil {
				return fmt.Errorf("failed to create job group %s: %w", jobGroup.Name, err)
			}
		}

		// Title seed data ekleme
		titles := []models.Title{
			// Doktor unvanları (JobGroupID: 1)
			{Name: "Başhekim", JobGroupID: 1},
			{Name: "Uzman Doktor", JobGroupID: 1},
			{Name: "Pratisyen Hekim", JobGroupID: 1},
			{Name: "Asistan Doktor", JobGroupID: 1},
			// Hemşire unvanları (JobGroupID: 2)
			{Name: "Başhemşire", JobGroupID: 2},
			{Name: "Sorumlu Hemşire", JobGroupID: 2},
			{Name: "Hemşire", JobGroupID: 2},
			// Teknisyen unvanları (JobGroupID: 3)
			{Name: "Laborant", JobGroupID: 3},
			{Name: "Radyoloji Teknisyeni", JobGroupID: 3},
			{Name: "Anestezi Teknisyeni", JobGroupID: 3},
			// İdari Personel unvanları (JobGroupID: 4)
			{Name: "İnsan Kaynakları Uzmanı", JobGroupID: 4},
			{Name: "Muhasebe Uzmanı", JobGroupID: 4},
			{Name: "Hasta Kabul", JobGroupID: 4},
			// Güvenlik unvanları (JobGroupID: 5)
			{Name: "Güvenlik Amiri", JobGroupID: 5},
			{Name: "Güvenlik Görevlisi", JobGroupID: 5},
		}

		for _, title := range titles {
			if err := db.Create(&title).Error; err != nil {
				return fmt.Errorf("failed to create title %s: %w", title.Name, err)
			}
		}

		fmt.Println("Personnel seed data created successfully!")
	}

	return nil
}
