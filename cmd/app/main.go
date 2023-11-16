package main

import (
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/ports/consumer"
	"FIO_App/pkg/ports/graph"
	"FIO_App/pkg/ports/rest"
	"FIO_App/pkg/repo"
	"FIO_App/pkg/service"
	"FIO_App/pkg/storage/database/postgres"
	"FIO_App/pkg/storage/database/redis"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger = logging.GetLogger()

func main() {
	ctx := context.Background()

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
	logger.Info("connected to topic FIO successfully")

	wg := &sync.WaitGroup{}

	wg.Add(3)

	go func() {
		defer wg.Done()
		//wait for kafka server to start
		time.Sleep(10 * time.Second)
		c.ConsumeMessages(ctx)
	}()

	r := rest.SetupRoutes(repository, pFio, logger)
	g, err := graph.NewGraphQLServer(fmt.Sprintf(":%s", os.Getenv("GRAPHQL_PORT")), repository, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// configuring graceful shutdown
	sigQuit := make(chan os.Signal, 1)
	defer close(sigQuit)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	go func() {
		defer wg.Done()
		logger.Info("started rest server successfully")
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		logger.Info("started graphql server successfully")
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("GRAPHQL_PORT")), g)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	wg.Wait()

}
