package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error){
	cli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	if cli.Ping(ctx).Err() == nil {
		return nil, errors.New("connect redis failed")
	}
	return cli, nil
}