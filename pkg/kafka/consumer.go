package kafka

import (
	"FIO_App/pkg/storage/person"
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"time"
)

type FIO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

const (
	FioTopic    = "FIO"
	FailedTopic = "FIO_FAILED"
)

func ConsumeMessage(brokerUrl []string, storage person.IStorage) {
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
				if err = storage.CreatePerson(personDTO); err != nil {
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
			case msg := <-consumer.Messages():
				log.Println(msg.Value)
			case <-signals:
				log.Println("Interrupted")
				return
			}
		}
	}()

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
