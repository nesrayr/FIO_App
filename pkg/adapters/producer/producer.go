package producer

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/logging"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type failedFio struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Error      string `json:"error"`
}

type IProducer interface {
	SendFailedMessage(ctx context.Context, fio dtos.FIO, errorMsg string) error
	SendMessage(ctx context.Context, fio dtos.FIO) error
	Close() error
}

type Producer struct {
	w *kafka.Writer
	l logging.Logger
}

func NewProducer(brokerUrl string, topic string, l logging.Logger) *Producer {
	return &Producer{w: &kafka.Writer{
		Addr:  kafka.TCP(brokerUrl),
		Topic: topic,
	},
		l: l}
}

func (p *Producer) SendFailedMessage(ctx context.Context, fio dtos.FIO, errorMsg string) error {
	fioToSend := failedFio{
		Name:       fio.Name,
		Surname:    fio.Surname,
		Patronymic: fio.Patronymic,
		Error:      errorMsg,
	}
	data, err := json.Marshal(fioToSend)
	if err != nil {
		p.l.Errorf("cannot send fio %s %s to FIO_FAILED: %s", fio.Name, fio.Surname, err.Error())
		return err
	}
	err = p.w.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		p.l.Errorf("cannot send fio %s %s to FIO_FAILED: %s", fio.Name, fio.Surname, err.Error())
		return err
	}
	p.l.Infof("send fio %s %s to FIO_FAILED", fio.Name, fio.Surname)
	return nil
}

func (p *Producer) SendMessage(ctx context.Context, fio dtos.FIO) error {
	data, err := json.Marshal(fio)
	if err != nil {
		p.l.Errorf("cannot send fio %s %s to FIO: %s", fio.Name, fio.Surname, err.Error())
		return err
	}
	err = p.w.WriteMessages(ctx, kafka.Message{Value: data})
	if err != nil {
		p.l.Errorf("cannot send fio %s %s to FIO: %s", fio.Name, fio.Surname, err.Error())
		return err
	}
	p.l.Infof("send fio %s %s to FIO", fio.Name, fio.Surname)
	return nil
}

func (p *Producer) Close() error {
	return p.w.Close()
}
