package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

var (
	configFile = "config/config.yml"
	configType = "yml"
)

type (
	Config struct {
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Database       Database `mapstructure:"database"`
		Clickhouse     Database `mapstructure:"clickhouse"`
		Log            Log      `mapstructure:"log"`
		Key            Key      `mapstructure:"key"`
		Cron           Cron     `mapstructure:"cron"`
		Redis          Redis    `mapstructure:"redis"`
		MongoDB        MongoDB  `mapstructure:"mongodb"`
		Prometheus     Http     `mapstructure:"prometheus"`
		Lark           Lark     `mapstructure:"lark"`
	}

	Server struct {
		Address string `mapstructure:"address"`
	}

	Database struct {
		Driver       string        `mapstructure:"driver"`
		Host         string        `mapstructure:"host"`
		Port         int           `mapstructure:"port"`
		User         string        `mapstructure:"username"`
		Password     string        `mapstructure:"password"`
		Database     string        `mapstructure:"database"`
		MaxIdleConns int           `mapstructure:"maxIdleConns"`
		MaxLeftTime  time.Duration `mapstructure:"maxLeftTime"`
		MaxOpenConns int           `mapstructure:"maxOpenConns"`
		Charset      string        `mapstructure:"charset"`
		TimeZone     string        `mapstructure:"timeZone"`
		Name         string        `mapstructure:"name"`
	}

	Log struct {
		FileName       string        `mapstructure:"fileName"`
		LogLevel       zapcore.Level `mapstructure:"logLevel"`
		MaxSizeMb      int           `mapstructure:"maxSizeMB"`
		MaxBackupCount int           `mapstructure:"maxBackupCount"`
		MaxKeepDays    int           `mapstructure:"maxKeepDays"`
	}

	Key struct {
		Type string `mapstructure:"type"`

		// ak & sk
		AccessKey string `mapstructure:"accessKey"`
		SecretKey string `mapstructure:"secretKey"`

		// basic auth
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	Cron struct {
		On bool `mapstructure:"on"`
	}

	Redis struct {
		PoolConfig `yaml:"pool" mapstructure:"pool"`

		Name         string        `yaml:"name" mapstructure:"name"` // redis name, for trace
		Proto        string        `yaml:"proto" mapstructure:"proto"`
		Addr         string        `yaml:"addr" mapstructure:"addr"`
		Auth         string        `yaml:"auth" mapstructure:"auth"`
		DialTimeout  time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
		ReadTimeout  time.Duration `yaml:"readTimeout" mapstructure:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout" mapstructure:"writeTimeout"`
		DB           int           `yaml:"db" mapstructure:"db"`
		SlowLog      time.Duration `yaml:"slowLog" mapstructure:"slowLog"`
	}

	// PoolConfig is the pool configuration struct.
	PoolConfig struct {
		// Active number of items allocated by the pool at a given time.
		// When zero, there is no limit on the number of items in the pool.
		Active int `yaml:"active" mapstructure:"active"`
		// Idle number of idle items in the pool.
		Idle int `yaml:"idle" mapstructure:"idle"`
		// If WaitTimeout is set and the pool is at the Active limit, then Get() waits WaitTimeout
		// until a item to be returned to the pool before returning.
		WaitTimeout time.Duration `yaml:"waitTimeout" mapstructure:"waitTimeout"`
		// If WaitTimeout is not set, then Wait effects.
		// if Wait is set true, then wait until ctx timeout, or default false and return directly.
		Wait bool `yaml:"wait" mapstructure:"wait"`
	}

	MongoDB struct {
		URI        string `yaml:"uri" mapstructure:"uri"`
		AuthSource string `yaml:"authSource" mapstructure:"authSource"`
		User       string `yaml:"user" mapstructure:"user"`
		Password   string `yaml:"password" mapstructure:"password"`
		Database   string `yaml:"database" mapstructure:"database"`

		MaxPoolSize uint64 `yaml:"maxPoolSize" mapstructure:"maxPoolSize"`
		MinPoolSize uint64 `yaml:"minPoolSize" mapstructure:"minPoolSize"`

		ConnectTimeoutMS int64 `yaml:"connectTimeoutMS" mapstructure:"connectTimeoutMS"`
		SocketTimeoutMS  int64 `yaml:"socketTimeoutMS" mapstructure:"socketTimeoutMS"`
	}

	Http struct {
		URL   string `yaml:"url" mapstructure:"url"`
		Token string `yaml:"token" mapstructure:"token"`
	}

	Etcd struct {
		Endpoints   []string      `yaml:"endpoints" mapstructure:"endpoints"`
		Username    string        `yaml:"username" mapstructure:"username"`
		Password    string        `yaml:"password" mapstructure:"password"`
		DialTimeout time.Duration `yaml:"dialTimeout" mapstructure:"dialTimeout"`
		Service     string        `yaml:"service" mapstructure:"service"`
	}

	Lark struct {
		AppID     string `yaml:"appID" mapstructure:"appID"`
		AppSecret string `yaml:"appSecret" mapstructure:"appSecret"`
	}
)

func NewConfig() Config {
	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable decode into config struct, %v", err)
	}
	return *conf
}

func InitConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func SetConfigFile(file string) {
	configFile = file
}
