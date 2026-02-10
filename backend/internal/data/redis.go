package data

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RDB *redis.Client

func InitRedis() error {
	addr := viper.GetString("redis.addr")
	if addr == "" {
		return fmt.Errorf("redis addr is missing")
	}

	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}
