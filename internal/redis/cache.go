package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// SetCache sets a value in Redis with an expiration time
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetCache retrieves a value from Redis
func GetCache(ctx context.Context, key string) (string, error) {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// GetListCache retrieves a list and count from Redis
func GetListCache(ctx context.Context, key string) (string, int, error) {
	// Retrieve value from Redis
	val, err := GetCache(ctx, key)
	if err == redis.Nil {
		return "", 0, nil
	} else if err != nil {
		return "", 0, err
	}

	// Retrieve count from Redis
	countStr, err := GetCache(ctx, fmt.Sprintf("cnt-%s", key))
	if err == redis.Nil {
		return "", 0, nil
	} else if err != nil {
		return "", 0, err
	}

	// Convert count to int
	count, convErr := strconv.Atoi(countStr)
	if convErr != nil {
		return "", 0, convErr
	}

	return val, count, nil
}

// DeleteCache deletes a single key from Redis
func DeleteCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// DeleteCacheByPrefix deletes all keys that match a given prefix
func DeleteCacheByPrefix(ctx context.Context, prefix string) error {
	iter := RedisClient.Scan(ctx, 0, fmt.Sprintf("%s*", prefix), 0).Iterator()
	for iter.Next(ctx) {
		if err := RedisClient.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// ClearAllCache clears the entire Redis database (use with caution!)
func ClearAllCache(ctx context.Context) error {
	return RedisClient.FlushDB(ctx).Err()
}
