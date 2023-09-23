package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Surname     string `gorm:"not null"`
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}
