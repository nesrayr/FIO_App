package graph

import (
	"FIO_App/pkg/logging"
	"FIO_App/pkg/repo"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func NewGraphQLServer(addr string, repository repo.IRepository, logger logging.Logger) (*handler.Handler, error) {
	query := rootQuery(repository, logger)
	mutation := rootMutation(repository, logger)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		return nil, err
	}

	h := handler.New(&handler.Config{Schema: &schema, Pretty: true})
	return h, nil
}
