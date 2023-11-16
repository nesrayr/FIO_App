package main

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/ports/consumer"
	"FIO_App/pkg/ports/graph"
	"FIO_App/pkg/repo"
	"FIO_App/pkg/router"
	"FIO_App/pkg/service"
	"FIO_App/pkg/storage/database/postgres"
	"FIO_App/pkg/storage/database/redis"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()

	logger := logging.GetLogger()

	postgresDB, err := postgres.ConnectDB()
	if err != nil {
		logger.Fatal("Failed to connect database\n", err)
	}

	redisDB := redis.NewStorage()
	if err = redisDB.Connect(); err != nil {
		logger.Fatal("Failed to connect Redis\n", err)
	}
	defer func() {
		_ = redisDB.Close()
	}()

	repository := repo.NewRepository(context.Background(), postgresDB, redisDB.Client, logger)

	pFio := producer.NewProducer(os.Getenv("ADDRESS"), "FIO", logger)
	defer func() {
		_ = pFio.Close()
	}()
	pFailed := producer.NewProducer(os.Getenv("ADDRESS"), "FIO_FAILED", logger)
	defer func() {
		_ = pFailed.Close()
	}()

	s := service.NewService(repository, pFailed)

	c := consumer.NewConsumer(os.Getenv("ADDRESS"), "FIO", s, logger)
	defer func() {
		_ = c.Close()
	}()

	wg := &sync.WaitGroup{}

	wg.Add(2)

	r := router.SetupRoutes(repository, pFio, logger)

	//go func() {
	//	defer wg.Done()
	//	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	//	if err != nil {
	//		logger.Fatal(err)
	//	}
	//}()

	go func() {
		defer wg.Done()
		//wait for kafka server to start
		time.Sleep(10 * time.Second)
		c.ConsumeMessages(ctx)
	}()

	graphqlServerAddr := fmt.Sprintf("%s:%s", os.Getenv("GRAPHQL_HOST"), os.Getenv("GRAPHQL_PORT"))
	graphqlServer, err := graph.NewGraphQLServer(graphqlServerAddr, repository, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		defer wg.Done()
		logger.Info("started graphql server successfully")
		err = graphqlServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	if err != nil {
		logger.Fatal(err)
	}

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)

	//err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	//if err != nil {
	//	logger.Fatal(err)
	//}
}
