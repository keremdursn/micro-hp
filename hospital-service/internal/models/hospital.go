package models

import "gorm.io/gorm"

type Hospital struct {
	gorm.Model
	Name        string `gorm:"not null"`
	TaxNumber   string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Phone       string `gorm:"unique;not null"`
	Address     string `gorm:"not null"`
	CityID      uint   `gorm:"not null"`
	City        City
	DistrictID  uint `gorm:"not null"`
	District    District
	Polyclinics []HospitalPolyclinic `gorm:"foreignKey:HospitalID"`
}
