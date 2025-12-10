package infrastructure

import (
	"context"
	"fmt"
	"test-go/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	Client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) (*RedisManager, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisManager{Client: client}, nil
}

// Set adds a key-value pair to Redis with a TTL.
func (r *RedisManager) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key.
func (r *RedisManager) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete removes a key from Redis.
func (r *RedisManager) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// Replace updates a key only if it exists (XX).
func (r *RedisManager) Replace(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.SetXX(ctx, key, value, ttl).Err()
}

// GetInfo retrieves Redis server info.
func (r *RedisManager) GetInfo(ctx context.Context) (string, error) {
	return r.Client.Info(ctx).Result()
}
