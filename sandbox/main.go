package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
	})

	ctx := context.Background()
	client.Set(ctx, "name", "Alice", time.Second * 30)
	result:= client.Get(ctx, "name").Val()

	fmt.Println(result)

	client.HMSet(ctx, "map", map[string]interface{}{"a": "b", "c": "d", "e": "f"})
	fmt.Println(client.HGetAll(ctx, "map").Val())
}
