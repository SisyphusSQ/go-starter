package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewCmd(ctx context.Context, args ...interface{}) *redis.Cmd {
	return redis.NewCmd(ctx, args)
}
