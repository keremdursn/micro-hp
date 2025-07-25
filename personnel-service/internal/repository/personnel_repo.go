package repository

import (
	"personnel-service/internal/database"
	"personnel-service/internal/dto"
	"personnel-service/internal/models"
	"gorm.io/gorm"
)

type PersonnelRepository interface {
	GetAllJobGroups() ([]models.JobGroup, error)
	GetAllTitlesByJobGroup(jobGroupID uint) ([]models.Title, error)

	IsTCOrPhoneExists(tc, phone string) (bool, error)
	GetJobGroupByID(id uint) (*models.JobGroup, error)
	GetTitleByID(id uint) (*models.Title, error)
	CountHospitalHeads(hospitalID uint) (int64, error)
	GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error)
	GetPolyclinicByID(id uint) (*models.Polyclinic, error)
	CreateStaff(staff *models.Staff) error

	GetStaffByID(id uint) (*models.Staff, error)
	IsTCOrPhoneExistsExcludeID(id uint, tc, phone string) (bool, error)
	UpdateStaff(staff *models.Staff) error
	DeleteStaff(staff *models.Staff) error

	ListStaffWithFilter(hospitalID uint, filter dto.StaffListFilter, page, size int) ([]models.Staff, error)
	CountStaffWithFilter(hospitalID uint, filter dto.StaffListFilter) (int64, error)
}

type personnelRepository struct {
	db *gorm.DB
}

func NewPersonnelRepository(db *gorm.DB) PersonnelRepository {
	return &personnelRepository{db: db}
}

func (r *personnelRepository) GetAllJobGroups() ([]models.JobGroup, error) {
	var groups []models.JobGroup
	if err := r.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *personnelRepository) GetAllTitlesByJobGroup(jobGroupID uint) ([]models.Title, error) {
	var titles []models.Title
	if err := r.db.Where("job_group_id = ?", jobGroupID).Find(&titles).Error; err != nil {
		return nil, err
	}
	return titles, nil
}

func (r *personnelRepository) IsTCOrPhoneExists(tc, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Staff{}).Where("tc = ? OR phone = ?", tc, phone).Count(&count).Error
	return count > 0, err
}

func (r *personnelRepository) GetJobGroupByID(id uint) (*models.JobGroup, error) {
	var JobGroup models.JobGroup
	if err := r.db.First(&JobGroup, id).Error; err != nil {
		return nil, err
	}
	return &JobGroup, nil
}

func (r *personnelRepository) GetTitleByID(id uint) (*models.Title, error) {
	var title models.Title
	if err := r.db.First(&title, id).Error; err != nil {
		return nil, err
	}
	return &title, nil
}

func (r *personnelRepository) CountHospitalHeads(hospitalID uint) (int64, error) {
	var count int64
	err := database.DB.Table("staffs").
		Joins("JOIN titles ON staffs.title_id = titles.id").
		Where("titles.name = ? AND staffs.hospital_id = ?", "Başhekim", hospitalID).
		Count(&count).Error
	return count, err
}

func (r *personnelRepository) GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error) {
	var hp models.HospitalPolyclinic
	if err := r.db.First(&hp, id).Error; err != nil {
		return nil, err
	}
	return &hp, nil
}

func (r *personnelRepository) GetPolyclinicByID(id uint) (*models.Polyclinic, error) {
	var p models.Polyclinic
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *personnelRepository) CreateStaff(staff *models.Staff) error {
	return r.db.Create(staff).Error
}

func (r *personnelRepository) GetStaffByID(id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.First(&staff, id).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *personnelRepository) IsTCOrPhoneExistsExcludeID(id uint, tc, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Staff{}).
		Where("id = ? AND (tc = ? OR phone = ?)", id, tc, phone).
		Count(&count).Error
	return count > 0, err
}

func (r *personnelRepository) UpdateStaff(staff *models.Staff) error {
	return r.db.Save(staff).Error
}

func (r *personnelRepository) DeleteStaff(staff *models.Staff) error {
	return r.db.Delete(staff).Error
}

// Filtreli ve sayfalı personel listesini getirir
func (r *personnelRepository) ListStaffWithFilter(hospitalID uint, filter dto.StaffListFilter, page, size int) ([]models.Staff, error) {
	query := r.db.Model(&models.Staff{}).Where("hospital_id = ?", hospitalID)

	// Filtreleme alanları
	if filter.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+filter.LastName+"%")
	}
	if filter.TC != "" {
		query = query.Where("tc ILIKE ?", "%"+filter.TC+"%")
	}
	if filter.JobGroupID != nil {
		query = query.Where("job_group_id = ?", *filter.JobGroupID)
	}
	if filter.TitleID != nil {
		query = query.Where("title_id = ?", *filter.TitleID)
	}

	var staffs []models.Staff
	// Sayfalama
	if err := query.Offset((page - 1) * size).Limit(size).Find(&staffs).Error; err != nil {
		return nil, err
	}
	return staffs, nil
}

// Aynı filtreyle toplam kayıt sayısını bulur
func (r *personnelRepository) CountStaffWithFilter(hospitalID uint, filter dto.StaffListFilter) (int64, error) {
	query := r.db.Model(&models.Staff{}).Where("hospital_id = ?", hospitalID)

	// Filtreleme alanları
	if filter.FirstName != "" {
		query = query.Where("first_name ILIKE ?", "%"+filter.FirstName+"%")
	}
	if filter.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+filter.LastName+"%")
	}
	if filter.TC != "" {
		query = query.Where("tc ILIKE ?", "%"+filter.TC+"%")
	}
	if filter.JobGroupID != nil {
		query = query.Where("job_group_id = ?", *filter.JobGroupID)
	}
	if filter.TitleID != nil {
		query = query.Where("title_id = ?", *filter.TitleID)
	}

	var totalCount int64
	err := query.Count(&totalCount).Error
	return totalCount, err
}
