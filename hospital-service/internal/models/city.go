package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name      string     `gorm:"unique;not null"`
	Districts []District `gorm:"foreignKey:CityID"`
} 