package main

import (
	"FIO_App/graph"
	"FIO_App/pkg/kafka"
	"FIO_App/pkg/router"
	"FIO_App/pkg/storage/database/postgres"
	"FIO_App/pkg/storage/database/redis"
	"FIO_App/pkg/storage/person"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
)

func main() {
	postgresDB, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect database\n", err)
	}

	redisDB := redis.NewStorage()
	if err = redisDB.Connect(); err != nil {
		log.Fatal("Failed to connect Redis\n", err)
	}

	st := person.NewStorage(postgresDB)

	go kafka.ConsumeMessage([]string{os.Getenv("ADDRESS")}, st)
	go kafka.ConsumeFailedMessage([]string{os.Getenv("ADDRESS")})

	r := router.SetupRoutes(st)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("after starting server")

}
