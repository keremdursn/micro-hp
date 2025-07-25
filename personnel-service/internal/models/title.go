package models

import "gorm.io/gorm"

type Title struct {
	gorm.Model
	Name       string  `gorm:"unique;not null"`
	JobGroupID uint    `gorm:"not null"`
	Staff      []Staff `gorm:"foreignKey:TitleID"`
}
