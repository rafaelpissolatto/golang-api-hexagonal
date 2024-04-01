package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// KeyCacheDuration expiration time for a key in cache
const KeyCacheDuration = time.Hour * 1

// RedisCache redis cache connection
type RedisCache struct {
	Client *redis.Client
}

// Set put a new key value pair in cache
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client.Set(ctx, key, value, expiration)
}

// Get returns the value of a key
func (r *RedisCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client.Get(ctx, key)
}
