package redis

import (
	"time"

	"github.com/redis/go-redis/v9"

	"go-starter/config"
)

// New 实例化新的redis v8
func New(config config.Config) *Client {
	conf := config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Auth,
		DB:           conf.DB,
		WriteTimeout: conf.WriteTimeout,
		ReadTimeout:  conf.ReadTimeout,
		MinIdleConns: conf.Idle,
		PoolSize:     conf.Active, //缩放连接数
		PoolTimeout:  time.Duration(conf.WaitTimeout),
		DialTimeout:  conf.DialTimeout,
	})
	rdb.PoolStats()
	return rdb
}
