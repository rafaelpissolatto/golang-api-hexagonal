package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisCacheMock redis cache mock
type RedisCacheMock struct{}

var (
	SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	GetFunc func(ctx context.Context, key string) *redis.StringCmd
)

// Set is the cache mock for Set func
func (rc *RedisCacheMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return SetFunc(ctx, key, value, expiration)
}

// Get is the cache mock for Get func
func (rc *RedisCacheMock) Get(ctx context.Context, key string) *redis.StringCmd {
	return GetFunc(ctx, key)
}
