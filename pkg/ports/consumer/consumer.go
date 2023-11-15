package consumer

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/service"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	r *kafka.Reader
	s service.IService
}

func NewConsumer(brokerUrl string, topic string, s service.IService) *Consumer {
	return &Consumer{r: kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerUrl},
		Topic:   topic,
	}),
		s: s}
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("disconnected from FIO")
			return
		default:
			msg, err := c.r.ReadMessage(ctx)
			if err != nil {
				log.Printf("cannot get message from FIO: %s", err.Error())
			}

			var fio dtos.FIO
			if err = json.Unmarshal(msg.Value, &fio); err != nil {
				log.Printf("cannot get message from FIO: %s", err.Error())
			}

			log.Printf("adding fio by FIO topic: %s %s", fio.Name, fio.Surname)

			err = c.s.FillFio(ctx, fio)
			if err != nil {
				log.Printf("something went wrong: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
