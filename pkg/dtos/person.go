package dtos

import (
	"FIO_App/pkg/models"
)

type FIO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type PersonDTO struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

func SetPersonDTO(name, surname, patronymic string, age int, gender, nationality string) PersonDTO {
	return PersonDTO{Name: name, Surname: surname, Patronymic: patronymic, Age: age, Gender: gender, Nationality: nationality}
}

func ToPerson(dto PersonDTO) models.Person {
	return models.Person{Name: dto.Name, Surname: dto.Surname, Patronymic: dto.Patronymic, Age: dto.Age, Gender: dto.Gender,
		Nationality: dto.Nationality}
}

func ToPersonDTO(person models.Person) PersonDTO {
	return PersonDTO{Name: person.Name, Surname: person.Surname, Patronymic: person.Patronymic, Age: person.Age,
		Gender: person.Gender, Nationality: person.Nationality}
}
