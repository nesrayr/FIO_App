package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"os"
	"time"
)

type IProducer interface {
	SendMessageToQueue(fioData FIO) error
}

type Producer struct {
	config *sarama.Config
}

func NewProducer(config *sarama.Config) *Producer {
	return &Producer{config: config}
}

func (p *Producer) SendMessageToQueue(fioData FIO) error {
	// wait for kafka server to start
	time.Sleep(time.Second * 10)

	brokerURL := os.Getenv("ADDRESS")
	topic := "FIO"
	config := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer([]string{brokerURL}, config)
	if err != nil {
		log.Println("error creating producer")
		return err
	}
	defer producer.Close()

	msgBytes, err := json.Marshal(fioData)
	if err != nil {
		return err
	}

	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgBytes),
	}
	log.Println("message was successfully sent")
	return nil
}

func SendToFailedQueue(producer sarama.AsyncProducer, msg []byte, errorMsg string, topic string) {
	failedMsg := map[string]interface{}{
		"original_message": string(msg),
		"error":            errorMsg,
	}
	messageBytes, err := json.Marshal(failedMsg)
	if err != nil {
		log.Printf("Error marshaling failed message: %v", err)
		return
	}
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
	}
}
