package services

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func GetRedisConnection(ctx context.Context, v viper.Viper) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", v.GetString("redis.host"), v.GetInt("redis.port")),
		Password: v.GetString("redis.password"),
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func CloseRedis(client *redis.Client) error {
	return client.Close()
}
