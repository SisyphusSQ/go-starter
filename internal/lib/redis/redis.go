package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/health"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Client struct{ *redis.Client }

func New(lifecycle fx.Lifecycle, c config.Config) (*Client, error) {
	conf := c.Redis
	result := &Client{}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			client := redis.NewClient(&redis.Options{Network: conf.Proto, Addr: conf.Addr, Password: conf.Auth, DB: conf.DB, WriteTimeout: conf.WriteTimeout, ReadTimeout: conf.ReadTimeout, MinIdleConns: conf.Idle, PoolSize: conf.Active, PoolTimeout: conf.WaitTimeout, DialTimeout: conf.DialTimeout, ContextTimeoutEnabled: true})
			if err := client.Ping(ctx).Err(); err != nil {
				_ = client.Close()
				return fmt.Errorf("ping redis: %w", err)
			}
			result.Client = client
			return nil
		},
		OnStop: func(context.Context) error {
			if result.Client == nil {
				return nil
			}
			return result.Close()
		},
	})
	return result, nil
}
func Readiness(c *Client) health.Check {
	return health.Check{Name: "redis", Enabled: true, Ping: func(ctx context.Context) error {
		if c.Client == nil {
			return errors.New("redis not started")
		}
		return c.Ping(ctx).Err()
	}}
}
