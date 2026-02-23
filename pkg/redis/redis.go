package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New(url string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: url,

	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis - New - ping: %w", err)
	}
	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
