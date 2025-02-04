package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/riskibarqy/go-template/config"
)

var ctx = context.Background()

// RedisClient holds the Redis client instance
var RedisClient *redis.Client

// Init initializes the Redis client
func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPassword,
	})

	// Test the connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalln(err)
	}
}
