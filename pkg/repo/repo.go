package repo

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/errs"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IRepository interface {
	AddPerson(personDTO dtos.PersonDTO) error
	GetPersonByID(ID int) (dtos.PersonDTO, error)
	GetPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error)
	UpdatePerson(ID int, personDTO dtos.PersonDTO) error
	DeletePerson(ID int) error
}

type Repository struct {
	permRepo
	cachedRepo
	ctx context.Context
}

func NewRepository(ctx context.Context, db *gorm.DB, rc *redis.Client) *Repository {
	return &Repository{
		permRepo:   &PermRepo{*db},
		cachedRepo: &CachedRepo{*rc},
		ctx:        ctx,
	}
}

func (r *Repository) AddPerson(personDTO dtos.PersonDTO) error {
	person, err := r.CreatePerson(personDTO)
	if err != nil {
		return err
	}

	err = r.SetPersonByKey(r.ctx, person)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPersonByID(ID int) (dtos.PersonDTO, error) {
	if person, err := r.GetPersonByKey(r.ctx, ID); err.Error() != errs.ErrorPersonNotFound && err != nil {
		return dtos.PersonDTO{}, err
	} else if err == nil {
		return person, nil
	}

	if person, err := r.SelectPersonByID(ID); err != nil {
		return dtos.PersonDTO{}, err
	} else {
		if err = r.SetPersonByKey(r.ctx, person); err != nil {
			return dtos.PersonDTO{}, err
		}
		return dtos.ToPersonDTO(person), nil
	}
}

func (r *Repository) GetPeople(limit, offset int, nationality, gender string) ([]dtos.PersonDTO, error) {
	return r.SelectPeople(limit, offset, nationality, gender)
}

func (r *Repository) UpdatePerson(ID int, personDTO dtos.PersonDTO) error {
	err := r.EditPerson(ID, personDTO)
	if err != nil {
		return err
	}

	person := dtos.ToPerson(personDTO)
	person.ID = ID
	err = r.SetPersonByKey(r.ctx, person)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeletePerson(ID int) error {
	err := r.RemovePerson(ID)
	if err != nil {
		return err
	}

	err = r.DeletePersonByKey(r.ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
