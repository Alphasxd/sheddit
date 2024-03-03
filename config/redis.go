package config

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var ctx = context.Background()

var RDB *redis.Client

func InitRedis(cfg *RedisConfig) {
	host := cfg.Host
	port := cfg.Port
	password := cfg.Password
	database := cfg.Database
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       database, // use default DB
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	zap.L().Info("redis init success")
}
