package repository

import (
	"auth-service/internal/models"
	"gorm.io/gorm"
)

type SubUserRepository interface {
	IsAuthorityExists(tc, email, phone string) (bool, error)
	CreateAuthority(authority *models.Authority) error
	GetByID(id uint) (*models.Authority, error)
	GetAllSubUsersByHospitalID(hospitalID uint) ([]models.Authority, error)
	IsUniqueForUpdate(id uint, tc, email, phone string) (bool, error)
	Update(user *models.Authority) error
	Delete(user *models.Authority) error
}

type subUserRepository struct {
	db *gorm.DB
}

func NewSubUserRepository(db *gorm.DB) SubUserRepository {
	return &subUserRepository{db: db}
}

func (r *subUserRepository) IsAuthorityExists(tc, email, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Authority{}).Where("tc = ? OR email = ? OR phone = ?", tc, email, phone).Count(&count).Error
	return count > 0, err
}

func (r *subUserRepository) CreateAuthority(authority *models.Authority) error {
	return r.db.Create(authority).Error
}

func (r *subUserRepository) GetByID(id uint) (*models.Authority, error) {
	var user models.Authority
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *subUserRepository) GetAllSubUsersByHospitalID(hospitalID uint) ([]models.Authority, error) {
	var users []models.Authority
	if err := r.db.Where("hospital_id = ?", hospitalID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *subUserRepository) IsUniqueForUpdate(id uint, tc, email, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Authority{}).
		Where("id != ? AND (tc = ? OR email != ? OR phone = ?)", id, tc, email, phone).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *subUserRepository) Update(user *models.Authority) error {
	return r.db.Save(user).Error
}

func (r *subUserRepository) Delete(user *models.Authority) error {
	return r.db.Delete(user).Error
}
