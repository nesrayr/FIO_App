package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"FIO_App/graph/model"
	"FIO_App/pkg/dtos"
	"context"
	"errors"
	"strconv"
)

// CreatePerson is the resolver for the createPerson field.
func (r *mutationResolver) CreatePerson(ctx context.Context, input model.PersonInput) (bool, error) {
	if input.Name == "" || input.Surname == "" {
		return false, errors.New("name and surname shouldn't be empty")
	}

	if err := r.storage.CreatePerson(dtos.SetPersonDTO(input.Name, input.Surname, input.Patronymic, input.Age,
		input.Gender, input.Nationality)); err != nil {
		return false, err
	}
	return true, nil
}

// DeletePerson is the resolver for the deletePerson field.
func (r *mutationResolver) DeletePerson(ctx context.Context, id string) (bool, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return false, errors.New("no such id")
	}

	if err = r.storage.DeletePerson(ID); err != nil {
		return false, err
	}

	return true, nil
}

// UpdatePerson is the resolver for the updatePerson field.
func (r *mutationResolver) UpdatePerson(ctx context.Context, id string, input model.PersonInput) (bool, error) {
	if input.Name == "" || input.Surname == "" {
		return false, errors.New("name and surname shouldn't be empty")
	}

	if err := r.storage.CreatePerson(dtos.SetPersonDTO(input.Name, input.Surname, input.Patronymic, input.Age,
		input.Gender, input.Nationality)); err != nil {
		return false, err
	}
	return true, nil
}

// GetPeople is the resolver for the getPeople field.
func (r *queryResolver) GetPeople(ctx context.Context, filter *model.PersonFilter, pagination *model.Pagination) ([]*model.Person, error) {
	var limit, offset int
	if pagination == nil {
		limit, offset = 0, 0
	} else {
		limit, offset = pagination.Limit, pagination.Offset
	}
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset should be greater than 0")
	}

	var gender, nationality string
	if filter == nil {
		gender, nationality = "", ""
	} else {
		gender, nationality = filter.Gender, filter.Nationality
	}

	people, err := r.storage.GetPeople(limit, offset, nationality, gender)
	if err != nil {
		return nil, err
	}
	var response []*model.Person
	for _, p := range people {
		cur := model.Person{Name: p.Name, Surname: p.Surname, Patronymic: p.Patronymic, Age: p.Age,
			Nationality: p.Nationality, Gender: p.Gender}
		response = append(response, &cur)
	}
	return response, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }