package consumer

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/logging"
	"FIO_App/pkg/service"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
	s service.IService
	l logging.Logger
}

func NewConsumer(brokerUrl string, topic string, s service.IService, l logging.Logger) *Consumer {
	return &Consumer{r: kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerUrl},
		Topic:   topic,
	}),
		s: s,
		l: l}
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.l.Error("disconnected from FIO")
			return
		default:
			msg, err := c.r.ReadMessage(ctx)
			if err != nil {
				c.l.Errorf("cannot get message from FIO: %s", err.Error())
			}

			var fio dtos.FIO
			if err = json.Unmarshal(msg.Value, &fio); err != nil {
				c.l.Errorf("cannot get message from FIO: %s", err.Error())
			}

			c.l.Infof("adding fio by FIO topic: %s %s", fio.Name, fio.Surname)

			err = c.s.FillFio(ctx, fio)
			if err != nil {
				c.l.Errorf("something went wrong: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
