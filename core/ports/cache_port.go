package ports

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// IRedis redis adapter
type IRedis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}
