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

func (r *RedisManager) GetStatus() map[string]interface{} {
	stats := make(map[string]interface{})
	if r == nil || r.Client == nil {
		stats["connected"] = false
		return stats
	}

	pong, err := r.Client.Ping(context.Background()).Result()
	stats["connected"] = err == nil
	stats["ping"] = pong
	stats["addr"] = r.Client.Options().Addr
	stats["db"] = r.Client.Options().DB

	// Add pool stats
	pool := r.Client.PoolStats()
	stats["pool_hits"] = pool.Hits
	stats["pool_misses"] = pool.Misses
	stats["pool_timeouts"] = pool.Timeouts
	stats["pool_total_conns"] = pool.TotalConns
	stats["pool_idle_conns"] = pool.IdleConns

	return stats
}

// GetInfo retrieves Redis server info.
func (r *RedisManager) GetInfo(ctx context.Context) (string, error) {
	return r.Client.Info(ctx).Result()
}

// ScanKeys returns a list of keys matching the pattern. Limit to 100 for safety.
func (r *RedisManager) ScanKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	iter := r.Client.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// GetValue returns the value of a specific key for monitoring.
// It assumes string for simplicity, but could be extended.
func (r *RedisManager) GetValue(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
