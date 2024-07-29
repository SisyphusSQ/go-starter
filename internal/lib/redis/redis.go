package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	*PoolConfig `yaml:"pool" json:"pool"`

	Name         string        `yaml:"name" json:"name"` // redis name, for trace
	Proto        string        `yaml:"proto" json:"proto"`
	Addr         string        `yaml:"addr" json:"addr"`
	Auth         string        `yaml:"auth" json:"auth"`
	DialTimeout  time.Duration `yaml:"dialTimeout" json:"dialTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
	DB           int           `yaml:"db" json:"db"`
	SlowLog      time.Duration `yaml:"slowLog" json:"slowLog"`
}

// PoolConfig is the pool configuration struct.
type PoolConfig struct {
	// Active number of items allocated by the pool at a given time.
	// When zero, there is no limit on the number of items in the pool.
	Active int `yaml:"active" json:"active"`
	// Idle number of idle items in the pool.
	Idle int `yaml:"idle" json:"idle"`
	// Close items after remaining item for this duration. If the value
	// is zero, then item items are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration `yaml:"idleTimeout" json:"idleTimeout"`
	// If WaitTimeout is set and the pool is at the Active limit, then Get() waits WatiTimeout
	// until a item to be returned to the pool before returning.
	WaitTimeout time.Duration `yaml:"waitTimeout" json:"waitTimeout"`
	// If WaitTimeout is not set, then Wait effects.
	// if Wait is set true, then wait until ctx timeout, or default flase and return directly.
	Wait bool `yaml:"wait" json:"wait"`
}

//New 实例化新的redis v8
func New(conf *Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Auth,
		DB:           conf.DB,
		WriteTimeout: conf.WriteTimeout,
		ReadTimeout:  conf.ReadTimeout,
		IdleTimeout:  time.Duration(conf.IdleTimeout),
		MinIdleConns: conf.Idle,
		PoolSize:     conf.Active, //缩放连接数
		PoolTimeout:  time.Duration(conf.WaitTimeout),
		DialTimeout:  conf.DialTimeout,
	})
	rdb.PoolStats()
	rdb.AddHook(&OpenTracingHook{
		cfg:    conf,
		status: rdb.PoolStats(),
	})
	return rdb
}
