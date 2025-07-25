package models

import "gorm.io/gorm"

type JobGroup struct {
	gorm.Model
	Name   string  `gorm:"unique;not null"`
	Titles []Title `gorm:"foreignKey:JobGroupID"`
	Staff  []Staff `gorm:"foreignKey:JobGroupID"`
}
