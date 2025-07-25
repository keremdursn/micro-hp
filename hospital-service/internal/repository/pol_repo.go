package repository

import (
	"errors"

	"hospital-service/internal/models"
	"gorm.io/gorm"
)

type PolyclinicRepository interface {
	GetAllPolyclinic() ([]models.Polyclinic, error)
	IsPolyclinicAlreadyAdded(hospitalID, PolyclinicID uint) (bool, error)
	GetPolyclinicByID(PolyclinicID uint) (*models.Polyclinic, error)
	CreateHospitalPolyclinic(hp *models.HospitalPolyclinic) error
	CountByHospitalID(hospitalID uint) (int64, error)
	GetPaginatedByHospitalID(hospitalID uint, page, size int) ([]models.HospitalPolyclinic, error)
	CountPersonnel(hospitalPolyclinicID uint) (int64, error)
	GetPersonnelGroupCounts(hospitalPolyclinicID uint) ([]struct {
		GroupName string
		Count     int
	}, error)
	GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error)
	Delete(hp *models.HospitalPolyclinic) error
}

type polyclinicRepository struct {
	db *gorm.DB
}

func NewPolyclinicRepository(db *gorm.DB) PolyclinicRepository {
	return &polyclinicRepository{db: db}
}

func (r *polyclinicRepository) GetAllPolyclinic() ([]models.Polyclinic, error) {
	var polys []models.Polyclinic
	if err := r.db.Find(&polys).Error; err != nil {
		return nil, err
	}
	return polys, nil
}

func (r *polyclinicRepository) IsPolyclinicAlreadyAdded(hospitalID, PolyclinicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.HospitalPolyclinic{}).
		Where("hospital_id = ? AND polyclinic_id = ?", hospitalID, PolyclinicID).
		Count(&count).Error
	return count > 0, err

}

func (r *polyclinicRepository) GetPolyclinicByID(PolyclinicID uint) (*models.Polyclinic, error) {
	var poly models.Polyclinic
	if err := r.db.First(&poly, PolyclinicID).Error; err != nil {
		return nil, errors.New("polyclinic not found")
	}
	return &poly, nil
}

func (r *polyclinicRepository) CreateHospitalPolyclinic(hp *models.HospitalPolyclinic) error {
	return r.db.Create(hp).Error
}

func (r *polyclinicRepository) CountByHospitalID(hospitalID uint) (int64, error) {
	var total int64
	err := r.db.Model(&models.HospitalPolyclinic{}).
		Where("hospital_id = ?", hospitalID).Count(&total).Error
	return total, err
}

// Belirli sayfa ve boyutta poliklinik kaydı çekiliyor
func (r *polyclinicRepository) GetPaginatedByHospitalID(hospitalID uint, page, size int) ([]models.HospitalPolyclinic, error) {
	var hps []models.HospitalPolyclinic
	err := r.db.Where("hospital_id = ?", hospitalID).
		Offset((page - 1) * size).Limit(size).
		Find(&hps).Error
	return hps, err
}

// Toplam personel sayısı
func (r *polyclinicRepository) CountPersonnel(hospitalPolyclinicID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Staff{}).
		Where("hospital_polyclinic_id = ?", hospitalPolyclinicID).
		Count(&count).Error

	return count, err
}

// Meslek grubu dağılımı
func (r *polyclinicRepository) GetPersonnelGroupCounts(hospitalPolyclinicID uint) ([]struct {
	GroupName string
	Count     int
}, error) {
	var groupCounts []struct {
		GroupName string
		Count     int
	}

	err := r.db.Table("staffs").
		Select("job_groups.name as group_name, COUNT(*) as count").
		Joins("JOIN job_groups ON staffs.job_group_id = job_groups.id").
		Where("staffs.hospital_polyclinic_id = ?", hospitalPolyclinicID).
		Group("job_groups.name").
		Scan(&groupCounts).Error

	return groupCounts, err
}

func (r *polyclinicRepository) GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error) {
	var hp models.HospitalPolyclinic
	if err := r.db.First(&hp, id).Error; err != nil {
		return nil, err
	}
	return &hp, nil
}

func (r *polyclinicRepository) Delete(hp *models.HospitalPolyclinic) error {
	return r.db.Delete(hp).Error
}
