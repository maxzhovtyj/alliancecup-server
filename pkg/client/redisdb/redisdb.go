package redisdb

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/zh0vtyj/allincecup-server/internal/config"
)

const (
	InventoryProducts = "InventoryProducts"
)

func NewClient(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	dsn := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ping := client.Ping(ctx)
	if err := ping.Err(); err != nil {
		return nil, err
	}

	return client, nil
}
