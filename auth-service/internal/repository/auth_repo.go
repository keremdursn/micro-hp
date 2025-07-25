package repository

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	IsAuthorityExists(tc, email, phone string) (bool, error)
	CreateAuthority(authority *models.Authority) error

	GetAuthorityByEmailOrPhone(credential string) (*models.Authority, error)

	GetAuthorityByPhone(phone string) (*models.Authority, error)
	UpdateAuthorityPassword(authority *models.Authority, hashedPassword string) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) IsAuthorityExists(tc, email, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Authority{}).
		Where("tc = ? OR email = ? OR phone = ?", tc, email, phone).
		Count(&count).Error
	return count > 0, err
}

func (r *authRepository) CreateAuthority(a *models.Authority) error {
	return r.db.Create(a).Error
}

func (r *authRepository) GetAuthorityByEmailOrPhone(credential string) (*models.Authority, error) {
	var authority models.Authority
	err := r.db.Where("email = ? OR phone = ?", credential, credential).First(&authority).Error
	if err != nil {
		return nil, err
	}
	return &authority, nil
}

func (r *authRepository) GetAuthorityByPhone(phone string) (*models.Authority, error) {
	var authority models.Authority
	if err := r.db.Where("phone = ?", phone).First(&authority).Error; err != nil {
		return nil, err
	}
	return &authority, nil
}

func (r *authRepository) UpdateAuthorityPassword(authority *models.Authority, hashedPassword string) error {
	return r.db.Model(authority).Update("password", hashedPassword).Error
}
