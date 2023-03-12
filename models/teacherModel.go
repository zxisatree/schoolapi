package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Email    string    `gorm:"unique;"`
	Students []Student `gorm:"many2many:teachers_students;constraint:OnUpdate:CASCADE;"`
}
