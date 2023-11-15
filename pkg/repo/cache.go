package repo

import (
	"FIO_App/pkg/dtos"
	"FIO_App/pkg/errs"
	"FIO_App/pkg/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const expiration = time.Minute * 30

type cachedRepo interface {
	GetPersonByKey(ctx context.Context, key int) (dtos.PersonDTO, error)
	SetPersonByKey(ctx context.Context, person models.Person) error
	DeletePersonByKey(ctx context.Context, key int) error
}

type CachedRepo struct {
	redis.Client
}

func (r *CachedRepo) GetPersonByKey(ctx context.Context, key int) (dtos.PersonDTO, error) {
	personData, err := r.Get(ctx, strconv.Itoa(key)).Result()
	if err == redis.Nil {
		log.Printf("cannot find person with id %d in cache: not exist\n", key)
		return dtos.PersonDTO{}, errors.New(errs.ErrorPersonNotFound)
	} else if err != nil {
		log.Printf("cannot find person with id %d in cache: %s\n", key, err.Error())
		return dtos.PersonDTO{}, err
	}

	var receivedPerson dtos.PersonDTO
	if err = json.Unmarshal([]byte(personData), &receivedPerson); err != nil {
		log.Printf("cannot find person with id %d in cache: %s\n", key, err.Error())
		return dtos.PersonDTO{}, err
	}

	log.Printf("find person with id %d in cache\n", key)

	return receivedPerson, nil
}

func (r *CachedRepo) SetPersonByKey(ctx context.Context, person models.Person) error {
	dto := dtos.ToPersonDTO(person)
	data, err := json.Marshal(dto)
	if err != nil {
		log.Printf("cannot marshal data: %v\n", err)
		return err
	}

	if err = r.Set(ctx, strconv.Itoa(person.ID), data, expiration).Err(); err != nil {
		log.Printf("cannot insert person with id %d in cache: %s\n", person.ID, err.Error())
		return err
	}

	log.Printf("insert person with id %d in cache", person.ID)
	return nil
}

func (r *CachedRepo) DeletePersonByKey(ctx context.Context, key int) error {
	_, err := r.Del(ctx, strconv.Itoa(key)).Result()
	if err == redis.Nil || err == nil {
		return nil
	} else {
		log.Printf("cannot delete fio with id %d: %s", key, err.Error())
		return err
	}
}
