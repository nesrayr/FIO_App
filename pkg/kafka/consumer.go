package kafka

import (
	"FIO_App/pkg/repo"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"time"
)

const (
	FioTopic    = "FIO"
	FailedTopic = "FIO_FAILED"
)

type FIO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func ConsumeMessage(brokerUrl []string, r repo.IRepository) {
	topic := FioTopic

	// wait for kafka server to start
	time.Sleep(time.Second * 10)

	worker, err := ConnectConsumer(brokerUrl)
	if err != nil {
		log.Println(err)
	}
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Topic %s consumer started", topic)
	defer consumer.Close()

	producer, err := ConnectProducer(brokerUrl)
	defer producer.Close()

	signals := make(chan struct{})

	go func() {
		for {
			select {
			case msg := <-consumer.Messages():
				var fioData FIO
				if err = json.Unmarshal(msg.Value, &fioData); err != nil {
					log.Printf("Error parsing JSON: %v", err)
					SendToFailedQueue(producer, msg.Value, "Invalid JSON format", "FIO_FAILED")
					continue
				}
				if fioData.Name == "" || fioData.Surname == "" {
					errMsg := "Missing required fields: name and/or surname"
					log.Println(errMsg)
					SendToFailedQueue(producer, msg.Value, errMsg, "FIO_FAILED")
					continue
				}
				log.Printf("data: %v", fioData)
				personDTO, err := EnrichData(fioData)
				if err != nil {
					log.Println(err)
				}
				log.Println(personDTO)
				if err = r.AddPerson(personDTO); err != nil {
					log.Println(err)
				}
			case <-signals:
				log.Println("Interrupted")
				return
			}
		}
	}()
	<-signals
}

func ConsumeFailedMessage(brokerUrl []string) {
	// wait for kafka server to start
	time.Sleep(time.Second * 10)

	topic := FailedTopic

	worker, err := ConnectConsumer(brokerUrl)
	if err != nil {
		log.Println(err)
	}
	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Topic %s consumer started", topic)
	defer consumer.Close()

	signals := make(chan struct{})

	go func() {
		for {
			select {
			case msg, ok := <-consumer.Messages():
				if !ok {
					log.Println("Channel closed")
					return
				}
				log.Println(string(msg.Value))
			case <-signals:
				log.Println("Interrupted")
				return
			}
		}
	}()
	<-signals
}

func ConnectConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ConnectProducer(brokerUrl []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	conn, err := sarama.NewAsyncProducer(brokerUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
