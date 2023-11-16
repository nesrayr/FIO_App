package graph

import (
	"FIO_App/pkg/errs"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/repo"
	"errors"
	"github.com/graphql-go/graphql"
)

func rootQuery(repo repo.IRepository, logger logging.Logger) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
		Fields: graphql.Fields{
			"getFioByID": &graphql.Field{
				Type: fioType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if id, ok := p.Args["id"].(int); ok {
						logger.Infof("getting fio with id %d by graphql server", id)
						return repo.GetPersonByID(id)
					}
					return nil, errors.New(errs.ErrorInvalidInput)
				},
			},
			"getFios": &graphql.Field{
				Type: graphql.NewList(fioType),
				Args: graphql.FieldConfigArgument{
					"offset":      &graphql.ArgumentConfig{Type: graphql.Int},
					"limit":       &graphql.ArgumentConfig{Type: graphql.Int},
					"nationality": &graphql.ArgumentConfig{Type: graphql.String},
					"gender":      &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					offset := p.Args["offset"].(int)
					limit := p.Args["limit"].(int)
					nationality := p.Args["nationality"].(string)
					gender := p.Args["gender"].(string)

					logger.Infof("getting fios by graphql server")
					return repo.GetPeople(limit, offset, nationality, gender)
				},
			},
		},
	})
}
