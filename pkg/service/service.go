package service

import (
	"FIO_App/pkg/adapters/apis"
	"FIO_App/pkg/adapters/producer"
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/repo"
	"FIO_App/pkg/service/utils"
	"context"
)

type IService interface {
	FillFio(ctx context.Context, fio dtos.FIO) error
}

type Service struct {
	repo.IRepository
	producer.IProducer
}

func NewService(repo repo.IRepository, producer producer.IProducer) *Service {
	return &Service{repo, producer}
}

func (s *Service) FillFio(ctx context.Context, fio dtos.FIO) error {
	if err := utils.ValidateFIO(fio); err != nil {
		if sendErr := s.SendFailedMessage(ctx, fio, err.Error()); sendErr != nil {
			return sendErr
		}
		return err
	}

	person, err := apis.EnrichData(fio)
	if err != nil {
		if sendErr := s.SendFailedMessage(ctx, fio, err.Error()); sendErr != nil {
			return sendErr
		}
		return err
	}

	err = s.AddPerson(person)
	if err != nil {
		return err
	}

	return nil
}
