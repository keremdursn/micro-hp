package database

import (
	"fmt"

	"hospital-service/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	err := db.AutoMigrate(
		&models.Hospital{},
		&models.City{},
		&models.District{},
		&models.Polyclinic{},
		&models.HospitalPolyclinic{},
	)
	if err != nil {
		return err
	}

	// Seed data ekleme
	return seedData(db)
}

func seedData(db *gorm.DB) error {
	// Şehir verilerini kontrol et, yoksa ekle
	var cityCount int64
	db.Model(&models.City{}).Count(&cityCount)

	if cityCount == 0 {
		cities := []models.City{
			{Name: "İstanbul"},
			{Name: "Ankara"},
			{Name: "İzmir"},
			{Name: "Bursa"},
			{Name: "Antalya"},
		}

		for _, city := range cities {
			if err := db.Create(&city).Error; err != nil {
				return fmt.Errorf("failed to create city %s: %w", city.Name, err)
			}
		}

		// İlçe verilerini ekle
		districts := []models.District{
			// İstanbul ilçeleri (city_id: 1)
			{Name: "Kadıköy", CityID: 1},
			{Name: "Beşiktaş", CityID: 1},
			{Name: "Şişli", CityID: 1},
			{Name: "Üsküdar", CityID: 1},
			// Ankara ilçeleri (city_id: 2)
			{Name: "Çankaya", CityID: 2},
			{Name: "Keçiören", CityID: 2},
			{Name: "Yenimahalle", CityID: 2},
			// İzmir ilçeleri (city_id: 3)
			{Name: "Konak", CityID: 3},
			{Name: "Karşıyaka", CityID: 3},
			{Name: "Bornova", CityID: 3},
			// Bursa ilçeleri (city_id: 4)
			{Name: "Osmangazi", CityID: 4},
			{Name: "Nilüfer", CityID: 4},
			// Antalya ilçeleri (city_id: 5)
			{Name: "Muratpaşa", CityID: 5},
			{Name: "Kepez", CityID: 5},
		}

		for _, district := range districts {
			if err := db.Create(&district).Error; err != nil {
				return fmt.Errorf("failed to create district %s: %w", district.Name, err)
			}
		}

		fmt.Println("Seed data created successfully!")
	}

	// Polyclinic seed data ekleme
	var polyclinicCount int64
	db.Model(&models.Polyclinic{}).Count(&polyclinicCount)

	if polyclinicCount == 0 {
		polyclinics := []models.Polyclinic{
			{Name: "Dahiliye"},
			{Name: "Kardiyoloji"},
			{Name: "Nöroloji"},
			{Name: "Ortopedi"},
			{Name: "Göz Hastalıkları"},
			{Name: "Kulak Burun Boğaz"},
			{Name: "Üroloji"},
			{Name: "Kadın Hastalıkları ve Doğum"},
		}

		for _, polyclinic := range polyclinics {
			if err := db.Create(&polyclinic).Error; err != nil {
				return fmt.Errorf("failed to create polyclinic %s: %w", polyclinic.Name, err)
			}
		}

		// İlk hastane için örnek hospital polyclinics
		hospitalPolyclinics := []models.HospitalPolyclinic{
			{HospitalID: 1, PolyclinicID: 1}, // Dahiliye
			{HospitalID: 1, PolyclinicID: 2}, // Kardiyoloji
			{HospitalID: 1, PolyclinicID: 3}, // Nöroloji
			{HospitalID: 1, PolyclinicID: 4}, // Ortopedi
		}

		for _, hp := range hospitalPolyclinics {
			if err := db.Create(&hp).Error; err != nil {
				return fmt.Errorf("failed to create hospital polyclinic: %w", err)
			}
		}

		fmt.Println("Polyclinic seed data created successfully!")
	}

	return nil
}
