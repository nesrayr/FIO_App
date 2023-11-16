package graph

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/repo"
	"github.com/graphql-go/graphql"
)

func rootMutation(repo repo.IRepository, logger logging.Logger) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "mutation",
		Fields: graphql.Fields{
			"addFio": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"name":        &graphql.ArgumentConfig{Type: graphql.String},
					"surname":     &graphql.ArgumentConfig{Type: graphql.String},
					"patronymic":  &graphql.ArgumentConfig{Type: graphql.String},
					"age":         &graphql.ArgumentConfig{Type: graphql.Int},
					"gender":      &graphql.ArgumentConfig{Type: graphql.String},
					"nationality": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var payload dtos.PersonDTO
					payload.Name = p.Args["name"].(string)
					payload.Surname = p.Args["surname"].(string)
					payload.Patronymic = p.Args["patronymic"].(string)
					payload.Age = p.Args["age"].(int)
					payload.Gender = p.Args["gender"].(string)
					payload.Nationality = p.Args["nationality"].(string)
					logger.Info("adding fio by graphql server: %s %s", payload.Name, payload.Surname)
					if err := repo.AddPerson(payload); err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"updateFio": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id":          &graphql.ArgumentConfig{Type: graphql.Int},
					"name":        &graphql.ArgumentConfig{Type: graphql.String},
					"surname":     &graphql.ArgumentConfig{Type: graphql.String},
					"patronymic":  &graphql.ArgumentConfig{Type: graphql.String},
					"age":         &graphql.ArgumentConfig{Type: graphql.Int},
					"gender":      &graphql.ArgumentConfig{Type: graphql.String},
					"nationality": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					var payload dtos.PersonDTO
					payload.Name = p.Args["name"].(string)
					payload.Surname = p.Args["surname"].(string)
					payload.Patronymic = p.Args["patronymic"].(string)
					payload.Age = p.Args["age"].(int)
					payload.Gender = p.Args["gender"].(string)
					payload.Nationality = p.Args["nationality"].(string)
					logger.Info("updating fio with id %d by graphql server", id)
					if err := repo.UpdatePerson(id, payload); err != nil {
						return false, err
					}
					return true, nil
				},
			},
			"deleteFio": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					logger.Info("deleting fio with id %d by graphql server", id)
					if err := repo.DeletePerson(id); err != nil {
						return false, err
					}
					return true, nil
				},
			},
		},
	})
}
