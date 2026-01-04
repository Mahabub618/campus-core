package database

import (
	"context"
	"fmt"
	"time"

	"campus-core/internal/config"
	"campus-core/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

// ConnectRedis establishes a connection to Redis
func ConnectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	RedisClient = client
	logger.Info("Redis connected successfully", zap.String("addr", cfg.GetRedisAddr()))

	return client, nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// SetWithExpiry sets a key-value pair with expiration
func SetWithExpiry(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// Exists checks if a key exists
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Increment increments a counter and returns the new value
func Increment(ctx context.Context, key string) (int64, error) {
	return RedisClient.Incr(ctx, key).Result()
}

// SetNX sets a key only if it doesn't exist (for locking)
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return RedisClient.SetNX(ctx, key, value, expiration).Result()
}
