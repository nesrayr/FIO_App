package graph

import "github.com/graphql-go/graphql"

var fioType = graphql.NewObject(graphql.ObjectConfig{
	Name: "fio",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"name":        &graphql.Field{Type: graphql.String},
		"surname":     &graphql.Field{Type: graphql.String},
		"patronymic":  &graphql.Field{Type: graphql.String},
		"age":         &graphql.Field{Type: graphql.Int},
		"gender":      &graphql.Field{Type: graphql.String},
		"nationality": &graphql.Field{Type: graphql.String},
	},
})
