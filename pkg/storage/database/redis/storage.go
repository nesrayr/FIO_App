package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

type Storage struct {
	client *redis.Client
	ctx    context.Context
}

func NewStorage() *Storage {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})
	return &Storage{client: client, ctx: context.Background()}
}

func (s *Storage) Connect() error {
	err := s.client.Ping(s.ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	err := s.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close connnection to Redis: %w", err)
	}
	return nil
}
