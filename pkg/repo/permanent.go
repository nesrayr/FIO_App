package repo

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/models"
	"gorm.io/gorm"
)

type permRepo interface {
	CreatePerson(personDTO dtos.PersonDTO) (models.Person, error)
	RemovePerson(ID int) error
	EditPerson(ID int, personDTO dtos.PersonDTO) error
	SelectPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error)
	SelectPersonByID(ID int) (models.Person, error)
}

type PermRepo struct {
	gorm.DB
}

func (r *PermRepo) CreatePerson(personDTO dtos.PersonDTO) (models.Person, error) {
	var person models.Person
	res := r.Where("name=? AND surname=?", personDTO.Name, personDTO.Surname).First(&person)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return models.Person{}, res.Error
	}
	if res.Error == gorm.ErrRecordNotFound {
		person = dtos.ToPerson(personDTO)
		if err := r.Create(&person).Error; err != nil {
			return models.Person{}, err
		}
	}
	return person, nil
}

func (r *PermRepo) RemovePerson(ID int) error {
	var person models.Person
	res := r.Where("id=?", ID).Delete(&person)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *PermRepo) EditPerson(ID int, personDTO dtos.PersonDTO) error {
	var person models.Person
	if err := r.Where("id=?", ID).First(&person).Error; err != nil {
		return err
	}
	person.Name = personDTO.Name
	person.Surname = personDTO.Surname
	person.Patronymic = personDTO.Patronymic
	person.Age = personDTO.Age
	person.Gender = personDTO.Gender
	person.Nationality = personDTO.Nationality

	if err := r.Save(&person).Error; err != nil {
		return err
	}
	return nil
}

func (r *PermRepo) SelectPersonByID(ID int) (models.Person, error) {
	var person models.Person
	if err := r.Where("id=?", ID).First(&person).Error; err != nil {
		return models.Person{}, err
	}
	return person, nil
}

func (r *PermRepo) SelectPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error) {
	// default
	if limit == 0 {
		limit = 5
	}

	var people []models.Person
	var tx *gorm.DB
	if nationality == "" && gender == "" {
		tx = r.Table("people").Limit(limit).Offset(offset).Find(&people)
	} else if nationality == "" {
		tx = r.Table("people").Where("gender=?", gender).Limit(limit).Offset(offset).Find(&people)
	} else if gender == "" {
		tx = r.Table("people").Where("nationality=?", nationality).Limit(limit).Offset(offset).Find(&people)
	}
	if nationality != "" && gender != "" {
		tx = r.Table("people").Where("nationality=? AND gender=?", nationality, gender).Limit(limit).
			Offset(offset).Find(&people)
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	var result []dtos.PersonDTO
	for _, p := range people {
		result = append(result, dtos.ToPersonDTO(p))
	}

	return result, nil
}
