package main

import (
	"FIO_App/pkg/kafka"
	"FIO_App/pkg/router"
	"FIO_App/pkg/storage/database/postgres"
	"FIO_App/pkg/storage/database/redis"
	"FIO_App/pkg/storage/person"
	"fmt"
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

	//log.Println("sending message")
	//err = kafka.SendMessageToQueue()
	//if err != nil {
	//	log.Fatal(err)
	//}

	st := person.NewStorage(postgresDB)

	go kafka.ConsumeMessage([]string{os.Getenv("ADDRESS")}, st)
	go kafka.ConsumeFailedMessage([]string{os.Getenv("ADDRESS")})

	r := router.SetupRoutes(st)
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("after starting server")

}
