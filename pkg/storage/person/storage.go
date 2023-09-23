package person

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/storage/models"
	"gorm.io/gorm"
)

type IStorage interface {
	CreatePerson(personDTO dtos.PersonDTO) error
	DeletePerson(ID int) error
	EditPerson(ID int, personDTO dtos.PersonDTO) error
	GetPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error)
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Init() *gorm.DB {
	return s.db
}

func (s *Storage) CreatePerson(personDTO dtos.PersonDTO) error {
	db := s.Init()

	var person models.Person
	res := db.Where("name=?", personDTO.Name).First(&person)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	if res.Error == gorm.ErrRecordNotFound {
		person = dtos.ToPerson(personDTO)
		if err := db.Create(&person).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) DeletePerson(ID int) error {
	db := s.Init()

	var person models.Person
	res := db.Where("id=?", ID).Delete(&person)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Storage) EditPerson(ID int, personDTO dtos.PersonDTO) error {
	db := s.Init()

	var person models.Person
	if err := db.Where("id=?", ID).First(&person).Error; err != nil {
		return err
	}
	person.Name = personDTO.Name
	person.Surname = personDTO.Surname
	person.Patronymic = personDTO.Patronymic
	person.Age = personDTO.Age
	person.Gender = personDTO.Gender
	person.Nationality = personDTO.Nationality

	if err := db.Save(&person).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error) {
	db := s.Init()

	// default
	if limit == 0 {
		limit = 5
	}

	var people []models.Person
	var tx *gorm.DB
	if nationality == "" && gender == "" {
		tx = db.Table("people").Limit(limit).Offset(offset).Find(&people)
	} else if nationality == "" {
		tx = db.Table("people").Where("gender=?", gender).Limit(limit).Offset(offset).Find(&people)
	} else if gender == "" {
		tx = db.Table("people").Where("nationality=?", nationality).Limit(limit).Offset(offset).Find(&people)
	}
	if nationality != "" && gender != "" {
		tx = db.Table("people").Where("nationality=? AND gender=?", nationality, gender).Limit(limit).
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
