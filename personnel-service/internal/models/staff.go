package models

import "gorm.io/gorm"

type Staff struct {
	gorm.Model
	FirstName            string `gorm:"not null"`
	LastName             string `gorm:"not null"`
	TC                   string `gorm:"unique;not null"`
	Phone                string `gorm:"unique;not null"`
	WorkingDays          string `gorm:"not null"` // "1,2,3,4,5"
	HospitalID           uint   `gorm:"not null"`
	JobGroupID           uint   `gorm:"not null"`
	TitleID              uint   `gorm:"not null"`
	HospitalPolyclinicID *uint
}
