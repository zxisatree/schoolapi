package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email     string `gorm:"unique;"`
	Suspended bool
	Teachers  []Teacher `gorm:"many2many:teachers_students;constraint:OnUpdate:CASCADE;"`
}
