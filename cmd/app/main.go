package main

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/ports/consumer"
	"FIO_App/pkg/repo"
	"FIO_App/pkg/router"
	"FIO_App/pkg/service"
	"FIO_App/pkg/storage/database/postgres"
	"FIO_App/pkg/storage/database/redis"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()

	postgresDB, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect database\n", err)
	}

	redisDB := redis.NewStorage()
	if err = redisDB.Connect(); err != nil {
		log.Fatal("Failed to connect Redis\n", err)
	}
	defer func() {
		_ = redisDB.Close()
	}()

	repository := repo.NewRepository(context.Background(), postgresDB, redisDB.Client)
	//producerP := kafka.NewProducer(sarama.NewConfig())

	pFio := producer.NewProducer(os.Getenv("ADDRESS"), "FIO")
	defer func() {
		_ = pFio.Close()
	}()
	pFailed := producer.NewProducer(os.Getenv("ADDRESS"), "FIO_FAILED")
	defer func() {
		_ = pFailed.Close()
	}()

	s := service.NewService(repository, pFailed)

	c := consumer.NewConsumer(os.Getenv("ADDRESS"), "FIO", s)
	defer func() {
		_ = c.Close()
	}()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	r := router.SetupRoutes(repository, pFio)

	//go func() {
	//	defer wg.Done()
	//	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	go func() {
		defer wg.Done()
		//wait for kafka server to start
		time.Sleep(10 * time.Second)
		c.ConsumeMessages(ctx)
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	if err != nil {
		log.Fatal(err)
	}

	//go kafka.ConsumeMessage([]string{os.Getenv("ADDRESS")}, repository)
	//go kafka.ConsumeFailedMessage([]string{os.Getenv("ADDRESS")})

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)

	//err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	//if err != nil {
	//	log.Fatal(err)
	//}
}
