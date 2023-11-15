package models

type Person struct {
	ID          int    `gorm:"primary key"`
	Name        string `gorm:"not null"`
	Surname     string `gorm:"not null"`
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}
