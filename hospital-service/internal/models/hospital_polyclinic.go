package models

import "gorm.io/gorm"

type HospitalPolyclinic struct {
	gorm.Model
	HospitalID   uint    `gorm:"not null"`
	PolyclinicID uint    `gorm:"not null"`
	Staff        []Staff `gorm:"foreignKey:HospitalPolyclinicID"`
}
