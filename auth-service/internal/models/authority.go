package models

import "gorm.io/gorm"

type Authority struct {
	gorm.Model
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	TC         string `gorm:"unique;not null"`
	Email      string `gorm:"unique;not null"`
	Phone      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Role       string `gorm:"not null;default:'yetkili'"` // yetkili, calisan
	HospitalID uint   `gorm:"not null"`
}
